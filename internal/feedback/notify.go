package feedback

import (
	"github.com/gen2brain/beeep"
)

func Notify(requiredVerbosity int, title, message, icon string) {
	if int(cutOff) >= requiredVerbosity {
		err := beeep.Alert(title, message, icon)
		if err != nil {
			HandleWErr("error while beeping: %w", err)
		}
	}
}
