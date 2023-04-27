package formatter_test

import (
	"testing"

	"github.com/dcilke/hz/pkg/formatter"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	testcases := map[string]struct {
		noColor bool
		msg     map[string]any
		expect  string
	}{
		"message-no-color": {true, map[string]any{"message": "hello"}, "hello"},
		"msg-no-color":     {true, map[string]any{"msg": "hello"}, "hello"},
		"message-color":    {false, map[string]any{"message": "hello"}, "hello"},
		"msg-color":        {false, map[string]any{"msg": "hello"}, "hello"},
		"diff":             {true, map[string]any{"message": "hello", "msg": "there"}, "message=hello msg=there"},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			f := formatter.NewMessage(tc.noColor, formatKey)
			str := f.Format(tc.msg)

			require.Equal(t, tc.expect, str)
		})
	}
}

func TestMessage_ExcludeKeys(t *testing.T) {
	f := formatter.NewMessage(false, nil)
	require.Equal(t, []string{formatter.KeyMessage, formatter.KeyMsg}, f.ExcludeKeys())
}
