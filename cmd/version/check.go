package version

import (
	"github.com/jsnfwlr/keyper-cli/internal/app"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(CheckCmd)
}

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "check {{.AppName}}'s version",
	Run:   CheckRun,
}

func CheckRun(cmd *cobra.Command, args []string) {
	feedback.Print(feedback.Required, false, app.OutputVersion())
}
