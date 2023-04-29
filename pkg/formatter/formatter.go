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

// Fielder formats a key value pair
type Fielder func(key string, value any) string

// Stringer stringifies a value
type Stringer func(any) string

func Key(color bool) Stringer {
	return func(i any) string {
		return Colorize(fmt.Sprintf("%v=", i), ColorCyan, color)
	}
}

func Map(fn Stringer) Fielder {
	return func(key string, value any) string {
		ret := fn(key)
		switch v := value.(type) {
		case string:
			if needsQuote(v) {
				ret += strconv.Quote(v)
			} else {
				ret += v
			}
		case json.Number:
			ret += v.String()
		default:
			b, err := json.Marshal(v)
			if err == nil {
				ret += string(b)
			} else if strings.HasPrefix(err.Error(), "json: unsupported value: encountered a cycle") {
				ret += "<cycle>"
			} else {
				ret += fmt.Sprintf("%v", v)
			}
		}
		return ret
	}
}
