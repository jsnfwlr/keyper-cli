package app

import (
	"cmp"
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/jsnfwlr/keyper-cli/internal/feedback/style"
)

//go:embed templates/*.tmpl
var templates embed.FS

const (
	useDefault = false
	writeFile  = false
)

func HelpTemplate(cmd *cobra.Command) string {
	b, _ := templates.ReadFile("templates/help.tmpl")

	return string(b)
}

func VersionTemplate(cmd *cobra.Command) string {
	b, _ := templates.ReadFile("templates/version.tmpl")

	return postProcess(string(b))
}

func UsageTemplate(cmd *cobra.Command) string {
	if useDefault {
		return cmd.UsageTemplate()
	}

	cobra.AddTemplateFuncs(templateFuncs())

	b, _ := templates.ReadFile("templates/usage.tmpl")

	if writeFile {
		err := os.WriteFile("usage.custom", b, 0o644)
		if err != nil {
			fmt.Println(err)
		}
	}

	return string(b)
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"postProcess": postProcess,
		"green":       func(s string) string { return style.Apply(s, style.GreenFG) },
		"sub":         func(int1, int2 int) int { return int1 - int2 },
		"add":         func(int1, int2 int) int { return int1 + int2 },
		"latestVersion": func() string {
			_, ver, _ := GetNewVersionDetails()
			return cmp.Or[string](ver, "unknown")
		},
		"releaseURL": func() string {
			return ReleaseURL
		},
		"canUpdate": func() bool {
			ua, _, _ := GetNewVersionDetails()

			return ua
		},
	}
}

func postProcess(text string) string {
	// _, ver, _ := GetNewVersionDetails()
	// text = strings.ReplaceAll(text, "{{.LatestVersion}}", cmp.Or[string](ver, "unknown"))
	text = strings.ReplaceAll(text, "{{.AppName}}", AppName)
	text = strings.ReplaceAll(text, "{{.RootCmd}}", RootCmd)

	return text
}

func PadRight(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}

func PadLeft(s string, width int) string {
	return fmt.Sprintf("%*s", width, s)
}
