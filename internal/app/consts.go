package app

import (
	"fmt"
	"time"
)

const (
	ReleaseURL    = "https://github.com/jsnfwlr/keyper-cli/releases"
	ReleaseAPI    = "https://api.github.com/repos/jsnfwlr/keyper-cli/releases/latest"
	ReleaseFilter = ".tag_name"
	Release
)

const (
	StatusError int    = 300
	AppID       string = "keyper-cli"
	AppName     string = "Keyper CLI"
	AppDesc     string = "Interact with your Keyper service from the command line"
	URL         string = "https://jsnfwlr.com/projects/keyper-cli/"
	RootCmd     string = "keyper"
	RootShort   string = "Command Line Interface for Keyper"

	IncreaseVerbosityFlagName  = "verbosity"
	IncreaseVerbosityFlagShort = "v"
	DecreaseVerbosityFlagName  = "quietness"
	DecreaseVerbosityFlagShort = "q"
	MaxVerbosityFlagName       = "debug"
)

var (
	buildVersion    = "v0.0.1-nightly.1"
	buildBranch     = "stable" // [ stable | beta | alpha | nightly ]
	timeoutDuration = 10 * time.Second
)

func RootLong() string {
	return fmt.Sprintf("%s - %s\n%s\nVisit %s for user documentation and support", AppName, RootShort, AppDesc, URL)
}
