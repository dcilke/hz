package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Formatter defines a formatter for a specific pin key
type Formatter interface {
	Format(map[string]any) string
	ExcludeKeys() []string
}

// Extractor extracts multiple values and formats them
type Extractor func(map[string]any, string) string

// Stringer stringifies a value
type Stringer func(any) string

func Key(noColor bool) Stringer {
	return func(i any) string {
		return Colorize(fmt.Sprintf("%v=", i), ColorCyan, noColor)
	}
}

func Map(fn Stringer) Extractor {
	return func(m map[string]any, k string) string {
		ret := fn(k)
		switch fValue := m[k].(type) {
		case string:
			if needsQuote(fValue) {
				ret += strconv.Quote(fValue)
			} else {
				ret += fValue
			}
		case json.Number:
			ret += fValue.String()
		default:
			b, err := json.Marshal(fValue)
			if err == nil {
				ret += string(b)
			} else if strings.HasPrefix(err.Error(), "json: unsupported value: encountered a cycle") {
				ret += "<cycle>"
			} else {
				ret += fmt.Sprintf("%v", fValue)
			}
		}
		return ret
	}
}
