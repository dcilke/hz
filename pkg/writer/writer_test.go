package writer_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dcilke/golden"
	"github.com/dcilke/hz/pkg/writer"
	"github.com/stretchr/testify/require"
)

type j = map[string]any
type a = []any

func TestConsole(t *testing.T) {
	testcases := map[string]any{
		"level":      j{"level": "info"},
		"log.level":  j{"log": j{"level": "info"}},
		"timestamp":  j{"timestamp": "2022-08-03T12:34:25.605701107Z"},
		"time":       j{"time": "2022-08-03T12:34:25.605701107Z"},
		"@timestamp": j{"@timestamp": "2022-08-03T12:34:25.605701107Z"},
		"message":    j{"message": "message"},
		"msg":        j{"msg": "message"},
		"error":      j{"error": "error"},
		"err":        j{"err": "error"},
		"default": j{
			"foo":       "bar",
			"level":     "info",
			"timestamp": "2022-08-03T12:34:25.605701107Z",
			"message":   "message",
			"error":     "error",
		},
		"duplicates": j{
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
		"conflicts": j{
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
		"array": a{
			j{"level": "info", "message": "message"},
			j{"foo": "bar"},
			"foo",
			5,
		},
		"array-nested": a{
			j{"level": "info", "message": "message"},
			a{
				j{"level": "info", "message": "message"},
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := writer.New(
				writer.WithOut(buf),
				writer.WithNoColor(),
			)
			b, err := json.Marshal(tc)
			require.NoError(t, err)
			o, err := w.Write(b)
			require.True(t, o > 0)
			require.NoError(t, err)
			golden.Assert(t, buf.Bytes())
		})
	}
}
