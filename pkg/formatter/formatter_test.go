package formatter_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

func formatKey(key any) string {
	return fmt.Sprintf("%s=", key)
}

func TestKey(t *testing.T) {
	testcases := map[string]struct {
		color  bool
		key    any
		expect string
	}{
		"no-color": {false, "key", "key="},
		"color":    {true, "key", "\x1b[36mkey=\x1b[0m"},
		"number":   {false, 42, "42="},
		"byte":     {false, byte(195), "195="},
		"rune":     {false, rune(195), "195="},
		"object":   {false, map[string]any{"foo": "bar"}, "map[foo:bar]="},
		"array":    {false, []string{"foo", "bar"}, "[foo bar]="},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.Key(tc.color)
			require.Equal(t, tc.expect, f(tc.key))
		})
	}
}

func TestMap(t *testing.T) {
	var (
		mapCycle   = make(map[string]any)
		sliceCycle = []any{nil}
	)
	mapCycle["x"] = mapCycle
	sliceCycle[0] = sliceCycle

	data := map[string]any{
		"key": "value",
		"map": map[string]any{
			"foo": "bar",
		},
		"quotes":      `this is a string`,
		"num":         42,
		"nan":         math.NaN(),
		"pinf":        math.Inf(1),
		"ninf":        math.Inf(-1),
		"map-cycle":   mapCycle,
		"slice-cycle": sliceCycle,
	}

	testcases := map[string]struct {
		key    string
		expect string
	}{
		"key":         {"key", "key=value"},
		"map":         {"map", "map={\"foo\":\"bar\"}"},
		"quotes":      {"quotes", "quotes=\"this is a string\""},
		"num":         {"num", "num=42"},
		"nan":         {"nan", "nan=NaN"},
		"pinf":        {"pinf", "pinf=+Inf"},
		"ninf":        {"ninf", "ninf=-Inf"},
		"map-cycle":   {"map-cycle", "map-cycle=<cycle>"},
		"slice-cycle": {"slice-cycle", "slice-cycle=<cycle>"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.Map(formatKey)
			require.Equal(t, tc.expect, f(tc.key, data[tc.key]))
		})
	}
}
