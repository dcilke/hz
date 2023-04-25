package writer

import (
	"fmt"
	"strconv"
)

const (
	KeyError = "error"
	KeyErr   = "err"
)

type errorFormatter struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func newErrorFormatter(noColor bool, formatKey Stringer) Formatter {
	return &errorFormatter{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyError, KeyErr},
	}
}

func (f *errorFormatter) Format(m map[string]any, _ string) string {
	var ferr string
	var ferror string
	if i, ok := m[KeyError]; ok {
		ferror = f.formatValue(i)
	}
	if i, ok := m[KeyErr]; ok {
		ferr = f.formatValue(i)
	}
	if ok, value := sameOrEmpty(ferror, ferr); ok {
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

func (f *errorFormatter) ExcludeKeys() []string {
	return f.keys
}

func (f *errorFormatter) formatValue(i any) string {
	str, err := strconv.Unquote(fmt.Sprintf("%s", i))
	if err != nil {
		return colorize(fmt.Sprintf("%s", i), ColorRed, f.noColor)
	}
	return colorize(str, ColorRed, f.noColor)
}
