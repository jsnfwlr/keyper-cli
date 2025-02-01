package keys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/files"
	"github.com/jsnfwlr/keyper-cli/internal/keyper"
	"github.com/jsnfwlr/keyper-cli/internal/keys"
	"github.com/jsnfwlr/keyper-cli/internal/prompter"

	"github.com/jsnfwlr/go-user"
	"github.com/spf13/cobra"
	enumFlag "github.com/thediveo/enumflag/v2"
)

var keyType keys.TypeId

func init() {
	defaultComment := ""
	defaultUsername := ""
	defaultFilename := "~/.ssh/id_<type>"

	// set the default comment based on the current user and hostname
	u, _ := user.Username()

	d, _ := user.HomeDir()

	h, _ := os.Hostname()

	if h != "" && u != "" {
		defaultComment = fmt.Sprintf("%s@%s", u, h)
	}

	if u != "" && d != "" {
		defaultUsername = u
		defaultFilename = filepath.Join(d, ".ssh", "id_<type>")
	}

	NewCmd.Flags().VarP(enumFlag.New(&keyType, "type", keys.Types(), enumFlag.EnumCaseInsensitive), "type", "t", "encryption algorithm to use for generating the key\n   ")
	NewCmd.Flags().StringP("filename", "f", defaultFilename, "specifies the filename of the key file\n   ")
	NewCmd.Flags().BoolP("overwrite", "o", false, "overwrite the key file if it already exists")
	NewCmd.Flags().BoolP("no-add", "n", false, "do not add the public key to your account on the Keyper server")
	NewCmd.Flags().StringP("passphrase", "P", "", "specifies the passphrase to use for the key file (use \"\" for no passphrase)\n    (if not set, {{.AppName}} will prompt for a passphrase)")
	NewCmd.Flags().StringP("comment", "c", defaultComment, "specifies the comment to use for the key file\n   ")
	NewCmd.Flags().StringP("user", "u", defaultUsername, "specifies th user on the Keyper server to add the key to\n   ")
	NewCmd.Flags().IntP("bits", "b", 2048, "specifies the number of bits in the key to create.")

	NewCmd.MarkFlagsMutuallyExclusive("user", "no-add")

	BaseCmd.AddCommand(NewCmd)
}

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "generate files for a new SSH key-pair and add the public key to your account in keyper",
	// Long:             app.RootLong(),
	// PersistentPreRun: RootPreRun,
	Run: NewRun,
}

func NewRun(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	file, keyType, user, passphrase, comment, bitSize, addToKeyper, err := parseFlags(cmd)
	feedback.HandleFatalErr(err)

	x := []string{"Generating"}

	if addToKeyper {
		x = append(x, "uploading")
	}

	fmt.Printf("%s public/private %s key pair.\n", strings.Join(x, " and "), keyType)

	pk, fp, err := keys.Generate(ctx, file, keyType, comment, passphrase, bitSize)
	feedback.HandleFatalErr(err)

	feedback.Print(feedback.Info, false, "Your identification has been saved in %s", file)
	feedback.Print(feedback.Info, false, "Your public key has been saved in %s.pub", file)
	feedback.Print(feedback.Info, false, "The key fingerprint is:\n%s %s", fp, comment)
	feedback.Print(feedback.Extra, false, "Public key:\n%s", strings.TrimSpace(pk))

	_, err = exec.LookPath("ssh-keygen")
	if err == nil {
		cmd := exec.Command("ssh-keygen", "-lvf", file)
		output := &bytes.Buffer{}

		cmd.Stdout = output
		cmd.Stderr = output

		err := cmd.Start()
		feedback.HandleFatalErr(err)

		err = cmd.Wait()
		feedback.HandleFatalErr(err)

		var lines string
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, ":") {
				continue
			}
			lines += line + "\n"
		}

		if err := scanner.Err(); err == nil {
			feedback.Print(feedback.Info, false, "The key's randomart image is:\n%s", lines)
		}

	}

	if !addToKeyper {
		return
	}

	kc, err := keyper.NewClient()
	feedback.HandleFatalErr(err)

	existing, err := listKeys(kc, user)
	feedback.HandleFatalErr(err)

	var nkk keyper.SSHPublicKey

	msgAdd := "New key with name %s added to keyper"

	for _, k := range existing {
		if k.Name == comment {
			feedback.Print(feedback.Required, true, "Found existing key with name %s in keyper", comment)
			err = kc.RevokeSSHKey(user, k)
			feedback.HandleFatalErr(err)
			feedback.Print(feedback.Required, true, "Revoked existing key with name %s in keyper", comment)

			nkk = k

			msgAdd = "Replaced existing key with name %s in keyper"

			break
		}
	}

	nkk.Name = comment
	nkk.Key = strings.TrimSpace(pk)
	nkk.Fingerprint = fp

	if len(nkk.HostGroups) == 0 {
		usr, err := kc.GetUser(user)
		feedback.HandleFatalErr(err)

		options := []string{}
		for _, hg := range usr.Groups {
			g := hg.Parse()
			options = append(options, g.CN)
		}

		slices.Sort(options)

		prompt := prompter.New()
		for {
			sel, err := prompt.Select("Select the host group for the key", false, options...)
			feedback.HandleFatalErr(err)

			for _, hg := range usr.Groups {
				g := hg.Parse()
				if g.CN == sel {
					nkk.HostGroups = append(nkk.HostGroups, hg)
					break
				}
			}

			if !prompt.Bool("Add another host group?", true) {
				break
			}

		}
	}

	err = kc.AddSSHKey(user, nkk)
	feedback.HandleFatalErr(err)

	feedback.Print(feedback.Required, false, msgAdd, comment)
}

func parseFlags(cmd *cobra.Command) (filePath, algorithm, userName, passphrase, keyComment string, bitSize int, addToKeyper bool, fault error) {
	keyType := cmd.Flags().Lookup("type").Value.String()
	if keyType == "" {
		return "", "", "", "", "", 0, false, errors.New("could not parse value for --type flag")
	}

	kd := keys.GetDefinitions(keyType)
	prompt := prompter.New()

	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return "", "", "", "", "", 0, false, fmt.Errorf("could parse value for --filename flag: %w", err)
	}

	if strings.HasPrefix(filename, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", "", "", "", 0, false, fmt.Errorf("could not determine the home directory for %s: %w", filename, err)
		}

		filename = home + filename[1:]
	}

	if strings.HasSuffix(filename, "<type>") {
		filename = strings.Replace(filename, "<type>", kd.Flag(), 1)
	}

	if files.Exists(filename, true) && !cmd.Flags().Changed("overwrite") {
		q := fmt.Sprintf("%s already exists. Overwrite?", filename)
		if !prompt.Bool(q, false) {
			return "", "", "", "", "", 0, false, errors.New("file already exists")
		}
	}

	noAdd, err := cmd.Flags().GetBool("no-add")
	if err != nil {
		return "", "", "", "", "", 0, false, fmt.Errorf("could not parse value for --no-add flag: %w", err)
	}

	phrase, err := getPassphrase(cmd, prompt, false)
	if err != nil {
		return "", "", "", "", "", 0, false, err
	}

	bitsize, err := getBitsize(cmd, prompt, keyType)
	if err != nil {
		return "", "", "", "", "", 0, false, err
	}

	comment, err := getComment(cmd)
	if err != nil {
		return "", "", "", "", "", 0, false, err
	}

	user, err := getUser(cmd)
	if err != nil {
		return "", "", "", "", "", 0, false, err
	}

	return filename, keyType, user, phrase, comment, bitsize, !noAdd, nil
}

func getPassphrase(cmd *cobra.Command, prompt *prompter.Prompt, retry bool) (passphrase string, fault error) {
	var phrase string
	var err error

	if !retry { // if this is the first attempt to get the passphrase, check if the user has provided it
		// value will be "" if not provided, or provided as empty string
		phrase, err = cmd.Flags().GetString("passphrase")
		if err != nil {
			return "", fmt.Errorf("could not parse value for --passphrase flag: %w", err)
		}
	}

	if retry || !cmd.Flags().Changed("passphrase") { // if the user has not provided a passphrase, prompt for it
		phrase, err = prompt.Text("Enter passphrase (empty for no passphrase)", "", true)
		if err != nil {
			return "", fmt.Errorf("could not parse value for passphrase from prompt: %w", err)
		}

		p2, err := prompt.Text("Enter same passphrase again", "", true)
		if err != nil {
			return "", fmt.Errorf("could not parse value for passphrase from prompt: %w", err)
		}

		if phrase != p2 {
			feedback.Print(feedback.Required, false, "Passphrases do not match. Try Again.")
			return getPassphrase(cmd, prompt, true) // the user has entered two different passphrases, prompt again
		}
	}

	return phrase, nil
}

func getUser(cmd *cobra.Command) (userName string, fault error) {
	// get the user provided username
	user, err := cmd.Flags().GetString("user")
	if err != nil {
		return "", fmt.Errorf("could not parse value for --user flag: %w", err)
	}

	return user, nil
}

func getComment(cmd *cobra.Command) (keyComment string, fault error) {
	// get the user provided comment
	comment, err := cmd.Flags().GetString("comment")
	if err != nil {
		return "", fmt.Errorf("could not parse value for --comment flag: %w", err)
	}

	return comment, nil
}

func getBitsize(cmd *cobra.Command, prompt *prompter.Prompt, keyType string) (bitSize int, fault error) {
	var err error
	bitsize := 0

	// set the default bit size based on the key type
	switch keyType {
	case "rsa":
		bitsize = 2048
	case "dsa":
		bitsize = 1024
	case "ecdsa", "ecdsa-sk":
		bitsize = 256
	case "ed25519", "ed25519-sk":
		bitsize = 256
	}

	// if the user has not provided a bit size, return the default
	if !cmd.Flags().Changed("bits") {
		return bitsize, nil
	}

	// get the user provoded bit size
	bitsize, err = cmd.Flags().GetInt("bits")
	if err != nil {
		return 0, fmt.Errorf("could not parse value for --bits flag: %w", err)
	}

	// validate the user provided bit size based on the key type
	switch keyType {
	case "rsa":
		if bitsize < 1024 {
			bitsize, err = prompt.Number("RSA key size must be at least 1024 bits", 1024, false)
			if err != nil {
				return 0, fmt.Errorf("could update key bit-size from prompt: %w", err)
			}
		}
		if bitsize < 1024 {
			return 0, errors.New("RSA key size must be at least 1024 bits")
		}
	case "dsa":
		if bitsize != 1024 {
			if !prompt.Bool("DSA key size must be 1024 bits - continue with 1024 as bit size?", true) {
				return 0, errors.New("DSA key size must be 1024 bits")
			}
		}
		bitsize = 1024
	case "ecdsa", "ecdsa-sk":
		if !slices.Contains([]int{256, 384, 521}, bitsize) {
			feedback.Print(feedback.Required, false, "%d is an invalid bit size for ECDSA key - must be one of 256, 384, or 521", bitsize)
			options := []string{"256", "384", "521"}
			sel, err := prompt.Select("Select the curve size for the ECDSA key", false, options...)
			if err != nil {
				return 0, fmt.Errorf("could not parse value for curve size from prompt: %w", err)
			}

			switch sel {
			case "256":
				bitsize = 256
			case "384":
				bitsize = 384
			case "521":
				bitsize = 521
			}
		}
	case "ed25519", "ed25519-sk":
		if !prompt.Bool("ED25519 key size can not be changed from 256 bits - continue with 256 as bit size?", true) {
			return 0, errors.New("ED25519 key size must be 256 bits")
		}
		// the bit size for Ed25519 keys, as they're always 256 bits, and isn't used as an input parameter,
		// but may need to be used in fingerprinting etc.
		bitsize = 256
	}

	return bitsize, nil
}
