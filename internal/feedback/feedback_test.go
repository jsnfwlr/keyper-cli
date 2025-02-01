package feedback_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsnfwlr/keyper-cli/internal/app"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
)

func TestExceedsThreshold(t *testing.T) {
	var level, suppress feedback.Level

	for level = 0; level <= feedback.LevelCap; level++ {
		flag := ""

		if level > 0 {
			flag = fmt.Sprintf("-%s ", strings.Repeat("q", int(level)))
		}

		for suppress = 0; suppress <= feedback.LevelCap; suppress++ {
			feedback.SuppressNoise(suppress)
			limit := feedback.GetCutOff()
			t.Run(fmt.Sprintf("suppress=%d limit=%d level=%d", suppress, limit, level), func(t *testing.T) {
				expect := level > feedback.GetCutOff()
				check := feedback.ExceedsLimit(level)
				t.Logf("%s %s (suppress=%d limit=%d level=%d) expect %t outcome %t", app.AppName, flag, suppress, feedback.GetCutOff(), level, expect, check)

				assert.Equal(t, expect, check)
			})
		}
	}
}
