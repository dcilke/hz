package formatter_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

const defaultTimeFormat = "15:04:05"
const ts = "2022-08-03T12:34:25.142900417Z"
const expect = "12:34:25"
const cexpect = "\x1b[90m" + expect + "\x1b[0m"

func TestTimestamp(t *testing.T) {
	testcases := map[string]struct {
		color  bool
		msg    map[string]any
		expect string
	}{
		"timestamp-no-color":  {false, map[string]any{"timestamp": ts}, expect},
		"@timestamp-no-color": {false, map[string]any{"@timestamp": ts}, expect},
		"time-no-color":       {false, map[string]any{"time": ts}, expect},
		"timestamp-color":     {true, map[string]any{"timestamp": ts}, cexpect},
		"@timestamp-color":    {true, map[string]any{"@timestamp": ts}, cexpect},
		"time-color":          {true, map[string]any{"time": ts}, cexpect},
		"unknown-time":        {false, map[string]any{"time": "unknown"}, "unknown"},
		"number-time":         {false, map[string]any{"time": jn(1111)}, "00:18:31"},
		"invalid-time":        {false, map[string]any{"time": jn("unknown")}, "unknown"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.NewTimestamp(tc.color, formatKey, defaultTimeFormat)
			str := f.Format(tc.msg)
			require.Equal(t, tc.expect, str)
		})
	}
}

func TestTimestamp_ExcludeKeys(t *testing.T) {
	f := formatter.NewTimestamp(false, nil, defaultTimeFormat)
	require.Equal(t, []string{formatter.KeyTimestamp, formatter.KeyAtTimestamp, formatter.KeyTime}, f.ExcludeKeys())
}

func jn(n any) json.Number {
	if i, ok := n.(int); ok {
		return json.Number(strconv.Itoa(i))
	}
	if s, ok := n.(string); ok {
		return json.Number(s)
	}
	return json.Number(fmt.Sprintf("%v", n))
}
