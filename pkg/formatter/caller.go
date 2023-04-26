package formatter

import (
	"os"
	"path/filepath"
)

const (
	KeyCaller = "caller"
)

var _ Formatter = (*Caller)(nil)

type Caller struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func NewCaller(noColor bool, formatKey Stringer) Formatter {
	return &Caller{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyCaller},
	}
}

func (f *Caller) Format(m map[string]any, _ string) string {
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
			c = Colorize(c, ColorBold, f.noColor) + Colorize(" >", ColorCyan, f.noColor)
		}
		return c
	}
	return ""
}

func (f *Caller) ExcludeKeys() []string {
	return f.keys
}
