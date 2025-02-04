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

var attributeNames = map[Attribute]string{
	Bold:            "bold",
	Faint:           "faint",
	Italic:          "italic",
	Underline:       "underline",
	BlackFG:         "black foreground",
	RedFG:           "red foreground",
	GreenFG:         "green foreground",
	YellowFG:        "yellow foreground",
	BlueFG:          "blue foreground",
	MagentaFG:       "magenta foreground",
	CyanFG:          "cyan foreground",
	WhiteFG:         "white foreground",
	BrightBlackFG:   "bright black foreground",
	BrightRedFG:     "bright red foreground",
	BrightGreenFG:   "bright green foreground",
	BrightYellowFG:  "bright yellow foreground",
	BrightBlueFG:    "bright blue foreground",
	BrightMagentaFG: "bright magenta foreground",
	BrightCyanFG:    "bright cyan foreground",
	BrightWhiteFG:   "bright white foreground",
	BlackBG:         "black background",
	RedBG:           "red background",
	GreenBG:         "green background",
	YellowBG:        "yellow background",
	BlueBG:          "blue background",
	MagentaBG:       "magenta background",
	CyanBG:          "cyan background",
	WhiteBG:         "white background",
	BrightBlackBG:   "bright black background",
	BrightRedBG:     "bright red background",
	BrightGreenBG:   "bright green background",
	BrightYellowBG:  "bright yellow background",
	BrightBlueBG:    "bright blue background",
	BrightMagentaBG: "bright magenta background",
	BrightCyanBG:    "bright cyan background",
	BrightWhiteBG:   "bright white background",
}

func (a Attribute) Name() string {
	if attributeNames[a] == "" {
		return "unknown style attribute"
	}
	return attributeNames[a]
}

func Apply(s string, attrs ...Attribute) string {
	if len(attrs) == 0 {
		return s
	}

	parts := applyColors(attrs...)

	parts = append(parts, applyDecorations(attrs...)...)

	return string(Escape) + strings.Join(parts, string(Delimiter)) + string(Close) + s + string(Reset)
}

func applyColors(attrs ...Attribute) (colors []string) {
	var fg Attribute
	var bg Attribute
	parts := []string{}

	for _, a := range attrs {
		if strings.Contains(a.Name(), "foreground") {
			fg = a
		}

		if strings.Contains(a.Name(), "background") {
			bg = a
		}
	}

	if fg != "" {
		parts = append(parts, string(fg))
	}

	if bg != "" {
		parts = append(parts, string(bg))
	}

	return parts
}

func applyDecorations(attrs ...Attribute) (decorations []string) {
	parts := []string{}

	bold := false
	faint := false
	italic := false
	underline := false

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
		}
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

	return parts
}
