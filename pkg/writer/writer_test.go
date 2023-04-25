package writer_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dcilke/hz/pkg/writer"
	"github.com/stretchr/testify/require"
)

type j = map[string]any

func TestConsole(t *testing.T) {
	testcases := map[string]struct {
		in  j
		out string
	}{
		"level":      {j{"level": "info"}, "<nil> INF"},
		"log.level":  {j{"log": j{"level": "info"}}, "<nil> INF"},
		"timestamp":  {j{"timestamp": "2022-08-03T12:34:25.605701107Z"}, "12:34:25"},
		"time":       {j{"time": "2022-08-03T12:34:25.605701107Z"}, "12:34:25"},
		"@timestamp": {j{"@timestamp": "2022-08-03T12:34:25.605701107Z"}, "12:34:25"},
		"message":    {j{"message": "message"}, "<nil> message"},
		"msg":        {j{"msg": "message"}, "<nil> message"},
		"error":      {j{"error": "error"}, "<nil> error=error"},
		"err":        {j{"err": "error"}, "<nil> error=error"},
		"default": {
			j{
				"foo":       "bar",
				"level":     "info",
				"timestamp": "2022-08-03T12:34:25.605701107Z",
				"message":   "message",
				"error":     "error",
			},
			"12:34:25 INF message error=error foo=bar",
		},
		"duplicates": {
			j{
				"level":      "info",
				"timestamp":  "2022-08-03T12:34:25.605701107Z",
				"@timestamp": "2022-08-03T12:34:25.605701107Z",
				"time":       "2022-08-03T12:34:25.605701107Z",
				"message":    "message",
				"msg":        "message",
				"error":      "error",
				"err":        "error",
				"log": j{
					"level": "info",
				},
				"foo": "bar",
			},
			"12:34:25 INF message error=error foo=bar",
		},
		"conflicts": {
			j{
				"level":      "info",
				"timestamp":  "2022-08-03T00:00:00.000000000Z",
				"@timestamp": "2022-08-03T11:11:11.000000000Z",
				"time":       "2022-08-03T22:22:22.000000000Z",
				"message":    "message",
				"msg":        "msg",
				"error":      "error",
				"err":        "err",
				"log": j{
					"level": "warn",
				},
				"foo": "bar",
			},
			"timestamp=00:00:00 @timestamp=11:11:11 time=22:22:22 level=INF log.level=WRN message=message msg=msg error=error err=err foo=bar",
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := writer.New(
				writer.WithOut(buf),
				writer.WithNoColor(),
			)
			b, err := json.Marshal(tc.in)
			require.NoError(t, err)
			o, err := w.Write(b)
			require.True(t, o > 0)
			require.NoError(t, err)
			require.Equal(t, tc.out, buf.String())
		})
	}
}
