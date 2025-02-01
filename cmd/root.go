package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jsnfwlr/keyper-cli/cmd/keys"
	"github.com/jsnfwlr/keyper-cli/cmd/version"

	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	// "github.com/jsnfwlr/keyper-cli/cmd/keygen"
	"github.com/jsnfwlr/keyper-cli/internal/app"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:              app.RootCmd,
	Short:            app.RootShort,
	Long:             app.RootLong(),
	PersistentPreRun: RootPreRun,
}

func RootPreRun(cmd *cobra.Command, args []string) {
	if strings.HasSuffix(cmd.Short, "^") {
		feedback.Print(feedback.Required, false, "This command is a work in progress and may not work as expected.")
	}
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	_ = godotenv.Load()

	cobra.OnInitialize(initTool)

	cobra.EnableCaseInsensitive = true

	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().SortFlags = false
	RootCmd.PersistentFlags().SortFlags = false

	RootCmd.AddCommand(keys.BaseCmd)
	RootCmd.AddCommand(version.BaseCmd)

	RootCmd.PersistentFlags().Bool(app.MaxVerbosityFlagName, false, "increase the verbosity of the output to include all messages, including debug")
	RootCmd.PersistentFlags().CountP(app.IncreaseVerbosityFlagName, app.IncreaseVerbosityFlagShort, fmt.Sprintf("increase the verbosity of the output\ncan not be used with --%[2]s\n-%[1]s adds 'Extra' messages\n-%[1]s%[1]s adds 'Extra' and 'Debug' messages", app.IncreaseVerbosityFlagShort, app.DecreaseVerbosityFlagName))
	RootCmd.PersistentFlags().CountP(app.DecreaseVerbosityFlagName, app.DecreaseVerbosityFlagShort, fmt.Sprintf("reduce the verbosity of the output\ncan not be used with --%[2]s\n-%[1]s removes 'Info' messages\n-%[1]s%[1]s removes 'Info' and 'Warning' messages\n-%[1]s%[1]s%[1]s removes 'Info', 'Warning', and 'Error' messages", app.DecreaseVerbosityFlagShort, app.IncreaseVerbosityFlagName))

	RootCmd.PersistentFlags().Lookup(app.MaxVerbosityFlagName).Hidden = true

	RootCmd.MarkFlagsMutuallyExclusive(app.DecreaseVerbosityFlagName, app.MaxVerbosityFlagName, app.IncreaseVerbosityFlagName)

	cmdGroups := []*cobra.Group{
		{
			ID:    "core",
			Title: "Core Commands",
		},
		{
			ID:    "extra",
			Title: "Extra Commands",
		},

		// &cobra.Group{
		// 	ID:    "short",
		// 	Title: "Shortcuts",
		// },
	}

	RootCmd.AddGroup(cmdGroups...)

	RootCmd.SetHelpTemplate(app.HelpTemplate())
	RootCmd.SetVersionTemplate(app.VersionTemplate())
	RootCmd.SetUsageTemplate(app.UsageTemplate(RootCmd))
}

// initTool reads in config file and ENV variables if set.
func initTool() {
	level := 3

	quietness, _ := RootCmd.PersistentFlags().GetCount(app.DecreaseVerbosityFlagName)
	if quietness > 0 {
		// fmt.Printf("decreasing level by quietness: %d\n", quietness)
		level = level - quietness
	}

	verbosity, _ := RootCmd.PersistentFlags().GetCount(app.IncreaseVerbosityFlagName)
	if verbosity > 0 {
		// fmt.Printf("increasing level by verbosity: %d\n", verbosity)
		level = level + verbosity
	}

	debug, _ := RootCmd.PersistentFlags().GetBool("debug")
	if debug {
		// // fmt.Printf("setting level to debug\n")
		level = 6
	}

	feedback.SetVerbosity(feedback.NoiseLevel(level))

	// fmt.Printf("quietness: %d\n", quietness)
	// fmt.Printf("verbosity: %d\n", verbosity)
	// fmt.Printf("debug: %t\n", debug)

	// fmt.Printf("resulting level: %d/%d = %s\n", level, feedback.GetCutOff(), feedback.GetCutOff().String())

	// os.Exit(0)

	// feedback.Print(feedback.Debug, false, "Debug verbosity level: %s", feedback.GetCutOff().String())
	// feedback.Print(feedback.Extra, false, "Extra verbosity level: %s", feedback.GetCutOff().String())
	// feedback.Print(feedback.Info, false, "Info verbosity level: %s", feedback.GetCutOff().String())
	// feedback.Print(feedback.Warning, false, "Warning verbosity level: %s", feedback.GetCutOff().String())
	// feedback.Print(feedback.Error, false, "Error verbosity level: %s", feedback.GetCutOff().String())
}
