package keys

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/feedback/style"
	"github.com/jsnfwlr/keyper-cli/internal/keyper"
	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(ListCmd)
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list the SSH keys for a user in keyper",
	Run:   ListRun,
}

func ListRun(cmd *cobra.Command, args []string) {
	kc, err := keyper.NewClient()
	feedback.HandleFatalErr(err)

	keys, err := listKeys(kc)
	feedback.HandleFatalErr(err)

	feedback.Print(feedback.Info, false, "SSH keys:")
	for _, k := range keys {
		if k.Local {
			feedback.Print(feedback.Info, false, "%[1]d: %[5]s%[2]s%[6]s - %[3]s - %[4]s", k.KeyId, k.Name, k.Fingerprint, k.Date.Format("2006-01-02 15:04:05"), style.CyanFG, style.Reset)
		} else {
			feedback.Print(feedback.Info, false, "%[1]d: %[2]s - %[3]s - %[5]s%[4]s%[6]s", k.KeyId, k.Name, k.Fingerprint, k.Date.Format("2006-01-02 15:04:05"), style.RedFG, style.Reset)
		}
	}
}

func listKeys(kc *keyper.Client) (keys []keyper.SSHPublicKey, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current user: %[1]w", err)
	}

	h, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get hostname: %[1]w", err)
	}

	k, err := kc.GetSSHKeys(u.Username, h)
	if err != nil {
		return nil, fmt.Errorf("could not get keys: %[1]w", err)
	}

	return k, nil
}
