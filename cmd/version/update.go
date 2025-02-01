package version

import (
	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(UpdateCmd)
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "download the latest version of {{.AppName}} from github",
	Run:   UpdateRun,
}

func UpdateRun(cmd *cobra.Command, args []string) {
}
