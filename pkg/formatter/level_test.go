package formatter_test

import (
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {

	testcases := map[string]struct {
		color  bool
		msg    map[string]any
		expect string
	}{
		"trace-no-color":   {false, ml("trace"), "TRC"},
		"debug-no-color":   {false, ml("debug"), "DBG"},
		"info-no-color":    {false, ml("info"), "INF"},
		"warn-no-color":    {false, ml("warn"), "WRN"},
		"error-no-color":   {false, ml("error"), "ERR"},
		"fatal-no-color":   {false, ml("fatal"), "FTL"},
		"panic-no-color":   {false, ml("panic"), "PNC"},
		"unknown-no-color": {false, ml("unknown"), "UNK"},
		"10-no-color":      {false, ml(10), "TRC"},
		"20-no-color":      {false, ml(20), "DBG"},
		"30-no-color":      {false, ml(30), "INF"},
		"40-no-color":      {false, ml(40), "WRN"},
		"50-no-color":      {false, ml(50), "ERR"},
		"60-no-color":      {false, ml(60), "FTL"},
		"100-no-color":     {false, ml(100), "PNC"},
		"2-no-color":       {false, ml(2), "2"},
		"trace-color":      {true, ml("trace"), "\x1b[35mTRC\x1b[0m"},
		"debug-color":      {true, ml("debug"), "\x1b[33mDBG\x1b[0m"},
		"info-color":       {true, ml("info"), "\x1b[32mINF\x1b[0m"},
		"warn-color":       {true, ml("warn"), "\x1b[31mWRN\x1b[0m"},
		"error-color":      {true, ml("error"), "\x1b[1m\x1b[31mERR\x1b[0m\x1b[0m"},
		"fatal-color":      {true, ml("fatal"), "\x1b[1m\x1b[31mFTL\x1b[0m\x1b[0m"},
		"panic-color":      {true, ml("panic"), "\x1b[1m\x1b[31mPNC\x1b[0m\x1b[0m"},
		"unknown-color":    {true, ml("unknown"), "\x1b[1mUNK\x1b[0m"},
		"10-color":         {true, ml(10), "\x1b[35mTRC\x1b[0m"},
		"20-color":         {true, ml(20), "\x1b[33mDBG\x1b[0m"},
		"30-color":         {true, ml(30), "\x1b[32mINF\x1b[0m"},
		"40-color":         {true, ml(40), "\x1b[31mWRN\x1b[0m"},
		"50-color":         {true, ml(50), "\x1b[1m\x1b[31mERR\x1b[0m\x1b[0m"},
		"60-color":         {true, ml(60), "\x1b[1m\x1b[31mFTL\x1b[0m\x1b[0m"},
		"100-color":        {true, ml(100), "\x1b[1m\x1b[31mPNC\x1b[0m\x1b[0m"},
		"2-color":          {true, ml(2), "\x1b[1m2\x1b[0m"},
		"diff":             {false, map[string]any{"level": "info", "log": map[string]any{"level": "debug"}}, "level=INF log.level=DBG"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.NewLevel(tc.color, formatKey)
			str := f.Format(tc.msg)
			require.Equal(t, tc.expect, str)
		})
	}
}

func TestLevel_ExcludeKeys(t *testing.T) {
	f := formatter.NewLevel(false, nil)
	require.Equal(t, []string{"level", "log.level"}, f.ExcludeKeys())
}

func ml(level any) map[string]any {
	if l, ok := level.(int); ok {
		level = jn(l)
	}
	return map[string]any{
		formatter.KeyLevel: level,
		formatter.KeyLog: map[string]any{
			formatter.KeyLevel: level,
		},
	}
}
