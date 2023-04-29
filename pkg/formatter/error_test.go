package formatter_test

import (
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	testcases := map[string]struct {
		color  bool
		msg    map[string]any
		expect string
	}{
		"error-no-color": {false, map[string]any{"error": "err"}, "error=err"},
		"err-no-color":   {false, map[string]any{"err": "err"}, "error=err"},
		"error-color":    {true, map[string]any{"error": "err"}, "error=\x1b[31merr\x1b[0m"},
		"err-color":      {true, map[string]any{"err": "err"}, "error=\x1b[31merr\x1b[0m"},
		"diff":           {false, map[string]any{"error": "error", "err": "err"}, "error=error err=err"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.NewError(tc.color, formatKey)
			str := f.Format(tc.msg)

			require.Equal(t, tc.expect, str)
		})
	}
}

func TestError_ExcludeKeys(t *testing.T) {
	f := formatter.NewError(false, nil)
	require.Equal(t, []string{formatter.KeyError, formatter.KeyErr}, f.ExcludeKeys())
}
