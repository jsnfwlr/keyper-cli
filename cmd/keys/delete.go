package keys

import (
	"encoding/json"
	"fmt"

	"github.com/jsnfwlr/go-user"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/keyper"

	"github.com/spf13/cobra"
)

func init() {
	defaultUsername := ""

	u, _ := user.Username()

	if u != "" {
		defaultUsername = u
	}

	DeleteCmd.Flags().IntP("key-id", "k", 0, "the key ID to delete")
	DeleteCmd.Flags().StringP("user", "u", defaultUsername, "specifies th user on the Keyper server to remove the key from\n   ")

	BaseCmd.AddCommand(DeleteCmd)
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a key with the provided KeyID for a user in keyper",
	Run:   DeleteRun,
}

func DeleteRun(cmd *cobra.Command, args []string) {
	keyId, err := cmd.Flags().GetInt("key-id")
	feedback.HandleFatalErr(err)

	if keyId == 0 && !cmd.Flags().Changed("key-id") {
		feedback.HandleFatalErr(fmt.Errorf("key ID is required"))
	}

	u, err := cmd.Flags().GetString("user")
	feedback.HandleFatalErr(err)

	kc, err := keyper.NewClient()
	feedback.HandleFatalErr(err)

	keys, err := listKeys(kc, u)
	feedback.HandleFatalErr(err)

	for _, k := range keys {
		if k.KeyId == keyId {
			fmt.Printf("deleting key %d\n", keyId)

			b, _ := json.MarshalIndent(k, "", "  ")

			fmt.Printf("%s\n", string(b))

			err = kc.RevokeSSHKey(u, k)
			feedback.HandleFatalErr(err)

			return
		}
	}

	feedback.HandleFatalErr(fmt.Errorf("key ID %d not found", keyId))
}
