package keys

import (
	"fmt"
	"os"
	"time"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/feedback/style"
	"github.com/jsnfwlr/keyper-cli/internal/keyper"

	"github.com/jsnfwlr/go-user"
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

func listKeys(kc *keyper.Client) (keys []keyper.SSHPublicKey, err error) {
	u, err := user.Username()
	if err != nil {
		return nil, fmt.Errorf("could not get username: %[1]w", err)
	}

	h, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get hostname: %[1]w", err)
	}

	k, err := kc.GetSSHKeys(u, h)
	if err != nil {
		return nil, fmt.Errorf("could not get keys: %[1]w", err)
	}

	return k, nil
}
