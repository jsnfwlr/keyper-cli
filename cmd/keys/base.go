package keys

import (
	"github.com/spf13/cobra"
)

func init() {
}

var BaseCmd = &cobra.Command{
	Use:     "keys",
	GroupID: "core",
	Short:   "manage SSH keys",
}
