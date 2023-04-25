package writer

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

func colorize(s any, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func boldrize(s any, c int, disabled bool) string {
	return colorize(colorize(s, c, disabled), ColorBold, disabled)
}

func includes[T comparable](slice []T, value T) bool {
	for _, e := range slice {
		if e == value {
			return true
		}
	}
	return false
}

func sameOrEmpty[T comparable](args ...T) (bool, T) {
	zero := *new(T)
	vals := filter(args, func(i T) bool {
		return i != zero
	})
	if len(vals) == 0 {
		return true, zero
	}
	// if all args are all the same, return the value
	isSame := true
	for i := 1; i < len(vals); i++ {
		if vals[i-1] != vals[i] {
			isSame = false
			break
		}
	}
	if isSame {
		return true, vals[0]
	}
	// if there are multiple values, return the zero value
	return false, zero
}

func filter[T comparable](slice []T, filter func(T) bool) []T {
	var result []T
	for _, e := range slice {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

func kvJoin(slice ...string) string {
	var pairs []string
	for i := 0; i < len(slice); i += 2 {
		pairs = append(pairs, slice[i]+slice[i+1])
	}
	return strings.Join(pairs, " ")
}
