package style

import "strings"

type (
	Directive string
	Component string
	Attribute string
)

const (
	Escape    Component = "\033["
	Close     Component = "m"
	Delimiter Component = ";"

	Reset   Directive = "\033[0m"
	Reverse Directive = "\033[7m"

	Bold            Attribute = "1"
	Faint           Attribute = "2"
	Italic          Attribute = "3"
	Underline       Attribute = "4"
	BlackFG         Attribute = "30"
	RedFG           Attribute = "31"
	GreenFG         Attribute = "32"
	YellowFG        Attribute = "33"
	BlueFG          Attribute = "34"
	MagentaFG       Attribute = "35"
	CyanFG          Attribute = "36"
	WhiteFG         Attribute = "37"
	BrightBlackFG   Attribute = "90"
	BrightRedFG     Attribute = "91"
	BrightGreenFG   Attribute = "92"
	BrightYellowFG  Attribute = "93"
	BrightBlueFG    Attribute = "94"
	BrightMagentaFG Attribute = "95"
	BrightCyanFG    Attribute = "96"
	BrightWhiteFG   Attribute = "97"
	BlackBG         Attribute = "40"
	RedBG           Attribute = "41"
	GreenBG         Attribute = "42"
	YellowBG        Attribute = "43"
	BlueBG          Attribute = "44"
	MagentaBG       Attribute = "45"
	CyanBG          Attribute = "46"
	WhiteBG         Attribute = "47"
	BrightBlackBG   Attribute = "100"
	BrightRedBG     Attribute = "101"
	BrightGreenBG   Attribute = "102"
	BrightYellowBG  Attribute = "103"
	BrightBlueBG    Attribute = "104"
	BrightMagentaBG Attribute = "105"
	BrightCyanBG    Attribute = "106"
	BrightWhiteBG   Attribute = "107"
)

func (a Attribute) Name() string {
	switch a {
	case Bold:
		return "bold"
	case Faint:
		return "faint"
	case Italic:
		return "italic"
	case Underline:
		return "underline"
	case BlackFG:
		return "black foreground"
	case RedFG:
		return "red foreground"
	case GreenFG:
		return "green foreground"
	case YellowFG:
		return "yellow foreground"
	case BlueFG:
		return "blue foreground"
	case MagentaFG:
		return "magenta foreground"
	case CyanFG:
		return "cyan foreground"
	case WhiteFG:
		return "white foreground"
	case BrightBlackFG:
		return "bright black foreground"
	case BrightRedFG:
		return "bright red foreground"
	case BrightGreenFG:
		return "bright green foreground"
	case BrightYellowFG:
		return "bright yellow foreground"
	case BrightBlueFG:
		return "bright blue foreground"
	case BrightMagentaFG:
		return "bright magenta foreground"
	case BrightCyanFG:
		return "bright cyan foreground"
	case BrightWhiteFG:
		return "bright white foreground"
	case BlackBG:
		return "black background"
	case RedBG:
		return "red background"
	case GreenBG:
		return "green background"
	case YellowBG:
		return "yellow background"
	case BlueBG:
		return "blue background"
	case MagentaBG:
		return "magenta background"
	case CyanBG:
		return "cyan background"
	case WhiteBG:
		return "white background"
	case BrightBlackBG:
		return "bright black background"
	case BrightRedBG:
		return "bright red background"
	case BrightGreenBG:
		return "bright green background"
	case BrightYellowBG:
		return "bright yellow background"
	case BrightBlueBG:
		return "bright blue background"
	case BrightMagentaBG:
		return "bright magenta background"
	case BrightCyanBG:
		return "bright cyan background"
	case BrightWhiteBG:
		return "bright white background"
	}

	return "unknown style attribute"
}

func Apply(s string, attrs ...Attribute) string {
	if len(attrs) == 0 {
		return s
	}

	bold := false
	faint := false
	italic := false
	underline := false

	var fg Attribute
	var bg Attribute

	for _, a := range attrs {
		switch a {
		case Bold:
			bold = true
		case Faint:
			faint = true
		case Italic:
			italic = true
		case Underline:
			underline = true
		case RedFG, GreenFG, YellowFG, BlueFG, MagentaFG, CyanFG:
			if fg == "" {
				fg = a
			}
		case RedBG, GreenBG, YellowBG, BlueBG, MagentaBG, CyanBG:
			if bg == "" {
				bg = a
			}
		}
	}

	parts := []string{}

	if fg != "" {
		parts = append(parts, string(fg))
	}

	if bg != "" {
		parts = append(parts, string(bg))
	}

	if bold {
		parts = append(parts, string(Bold))
	}

	if faint {
		parts = append(parts, string(Faint))
	}

	if italic {
		parts = append(parts, string(Italic))
	}

	if underline {
		parts = append(parts, string(Underline))
	}

	return string(Escape) + strings.Join(parts, string(Delimiter)) + string(Close) + s + string(Reset)
}
