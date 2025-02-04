package keys

import (
	"github.com/spf13/cobra"
)

func init() {
	ArtCmd.Flags().StringP("filename", "f", "~/.ssh/id_<type>", "specifies the filename of the public key")

	BaseCmd.AddCommand(ArtCmd)
}

var ArtCmd = &cobra.Command{
	Use:    "art",
	Hidden: true,
	Short:  "show the random art from the public key provided",
	Run:    ArtRun,
}

func ArtRun(cmd *cobra.Command, args []string) {
	// pKey, _, err := ed25519.GenerateKey(nil)
	// feedback.HandleFatalErr(err)

	// pubSSH, err := ssh.NewPublicKey(pKey)
	// feedback.HandleFatalErr(err)

	// 	f, err := cmd.Flags().GetString("filename")
	// 	feedback.HandleFatalErr(err)

	// 	b, err := files.Read(f, true)
	// 	feedback.HandleFatalErr(err)

	// 	pk, _, _, _, err := ssh.ParseAuthorizedKey(b)
	// 	feedback.HandleFatalErr(err)

	// 	ppk, err := ssh.ParsePublicKey(pk.Marshal())
	// 	feedback.HandleFatalErr(err)

	// 	fp := ssh.FingerprintSHA256(ppk)

	// 	art := randomart.FromString(fp)

	// 	feedback.Print(feedback.Info, false, "%s", art)

	// 	pre := randomart.Transform("", "SHA256", 256, []byte(fp))

	// feedback.Print(feedback.Info, false, "%s", pre)
}
