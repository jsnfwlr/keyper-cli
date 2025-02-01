package keys

import (
	"fmt"
	"os"
	"time"

	"github.com/jsnfwlr/go-user"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/feedback/style"
	"github.com/jsnfwlr/keyper-cli/internal/keyper"

	"github.com/spf13/cobra"
)

func init() {
	defaultUsername := ""

	u, _ := user.Username()

	if u != "" {
		defaultUsername = u
	}

	ListCmd.Flags().StringP("user", "u", defaultUsername, "specifies th user on the Keyper server to add the key to\n   ")

	BaseCmd.AddCommand(ListCmd)
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list the SSH keys for a user in keyper",
	Run:   ListRun,
}

func ListRun(cmd *cobra.Command, args []string) {
	user, err := cmd.Flags().GetString("user")
	feedback.HandleFatalErr(err)

	kc, err := keyper.NewClient()
	feedback.HandleFatalErr(err)

	keys, err := listKeys(kc, user)
	feedback.HandleFatalErr(err)

	feedback.Print(feedback.Info, false, "SSH keys:")
	for _, k := range keys {
		name := k.Name
		if k.Local {
			name = style.Apply(name, style.CyanFG)
		}
		date := k.Date.Format("2006-01-02 15:04:05")
		if k.Date.Before(time.Now()) {
			date = style.Apply(date, style.RedFG)
		}
		feedback.Print(feedback.Info, false, "%d: %s - %s - %s", k.KeyId, name, k.Fingerprint, date)
	}
}

func listKeys(kc *keyper.Client, user string) (keys []keyper.SSHPublicKey, err error) {
	h, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get hostname: %[1]w", err)
	}

	k, err := kc.GetSSHKeys(user, h)
	if err != nil {
		return nil, fmt.Errorf("could not get keys: %[1]w", err)
	}

	return k, nil
}
