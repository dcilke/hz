package formatter_test

import (
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

func TestCaller(t *testing.T) {
	testcases := map[string]struct {
		noColor bool
		expect  string
	}{
		"no-color": {false, "caller >"},
		"color":    {true, "\x1b[1mcaller\x1b[0m\x1b[36m >\x1b[0m"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.NewCaller(tc.noColor)
			str := f.Format(map[string]any{
				formatter.KeyCaller: "caller",
				"extra":             "should be ignored",
			})

			require.Equal(t, tc.expect, str)
		})
	}
}

func TestCaller_ExcludeKeys(t *testing.T) {
	f := formatter.NewCaller(false)
	require.Equal(t, []string{formatter.KeyCaller}, f.ExcludeKeys())
}
