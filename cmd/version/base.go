package version

import (
	"github.com/spf13/cobra"
)

func init() {
}

var BaseCmd = &cobra.Command{
	Use:     "version",
	GroupID: "extra",
	Short:   "check {{.AppName}}'s version",
	Run:     CheckRun,
}
