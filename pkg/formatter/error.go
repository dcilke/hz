package formatter

import (
	"fmt"
	"strconv"

	"github.com/dcilke/gu"
)

const (
	KeyError = "error"
	KeyErr   = "err"
)

var _ Formatter = (*Error)(nil)

type Error struct {
	color     bool
	formatKey Stringer
	keys      []string
}

func NewError(color bool, formatKey Stringer) Formatter {
	return &Error{
		color:     color,
		formatKey: formatKey,
		keys:      []string{KeyError, KeyErr},
	}
}

func (f *Error) Format(m map[string]any) string {
	var ferr string
	var ferror string
	if i, ok := m[KeyError]; ok {
		ferror = f.formatValue(i)
	}
	if i, ok := m[KeyErr]; ok {
		ferr = f.formatValue(i)
	}
	if ok, value := gu.SameOrZero(ferror, ferr); ok {
		if value == "" {
			return ""
		}
		return f.formatKey(KeyError) + value
	}

	return kvJoin(
		f.formatKey(KeyError), ferror,
		f.formatKey(KeyErr), ferr,
	)
}

func (f *Error) ExcludeKeys() []string {
	return f.keys
}

func (f *Error) formatValue(i any) string {
	str, err := strconv.Unquote(fmt.Sprintf("%s", i))
	if err != nil {
		return Colorize(fmt.Sprintf("%s", i), ColorRed, f.color)
	}
	return Colorize(str, ColorRed, f.color)
}
