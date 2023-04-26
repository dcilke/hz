package formatter

import (
	"fmt"
	"strings"
)

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBold     = 1
	ColorDarkGray = 90
)

func Colorize(s any, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func Boldrize(s any, c int, disabled bool) string {
	return Colorize(Colorize(s, c, disabled), ColorBold, disabled)
}

func kvJoin(slice ...string) string {
	var pairs []string
	for i := 0; i < len(slice); i += 2 {
		pairs = append(pairs, slice[i]+slice[i+1])
	}
	return strings.Join(pairs, string(' '))
}

// needsQuote returns true when the string s should be quoted in output.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}
