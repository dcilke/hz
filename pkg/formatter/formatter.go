package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Formatter defines a formatter for a specific pin key
type Formatter interface {
	Format(map[string]any, string) string
	ExcludeKeys() []string
}

// Extractor extracts multiple values and formats them
type Extractor func(map[string]any, string) string

// Stringer stringifies a value
type Stringer func(any) string

func Key(noColor bool) Stringer {
	return func(i any) string {
		return Colorize(fmt.Sprintf("%s=", i), ColorCyan, noColor)
	}
}

func Map(noColor bool, fn Stringer) Extractor {
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
			if err != nil {
				ret += fmt.Sprintf(Colorize("[error: %v]", ColorRed, noColor), err)
			} else {
				ret += string(b)
			}
		}
		return ret
	}
}
