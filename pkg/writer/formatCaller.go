package writer

import (
	"os"
	"path/filepath"
)

const (
	KeyCaller = "caller"
)

type callerFormatter struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func newCallerFormatter(noColor bool, formatKey Stringer) Formatter {
	return &callerFormatter{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyCaller},
	}
}

func (f *callerFormatter) Format(m map[string]any, _ string) string {
	if i, ok := m[KeyCaller]; ok {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
			c = colorize(c, ColorBold, f.noColor) + colorize(" >", ColorCyan, f.noColor)
		}
		return c
	}
	return ""
}

func (f *callerFormatter) ExcludeKeys() []string {
	return f.keys
}
