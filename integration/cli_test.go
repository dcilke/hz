package integration

import (
	"testing"

	"github.com/dcilke/golden"
	"github.com/stretchr/testify/require"
)

var filecases = []string{"strings", "ndjson", "pretty", "array", "mixed"}
var levels = []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}

func TestCLI(t *testing.T) {
	for _, tc := range filecases {
		t.Run(tc, func(t *testing.T) {
			output, err := hz(fn(tc))
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Strict(t *testing.T) {
	for _, tc := range filecases {
		t.Run(tc, func(t *testing.T) {
			output, err := hz(fn(tc), "--strict")
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Level(t *testing.T) {
	for _, tc := range filecases {
		for _, level := range levels {
			t.Run(tc+"/"+level, func(t *testing.T) {
				output, err := hz(fn(tc), "--level", level)
				require.NoError(t, err)
				golden.Assert(t, output)
			})
		}
	}
}

func TestCLI_Help(t *testing.T) {
	testcases := []string{"--help", "-h"}
	for _, tc := range testcases {
		t.Run(tc, func(t *testing.T) {
			output, err := hz(tc)
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Flat(t *testing.T) {
	output, err := hz(fn("nested"), "--flat")
	require.NoError(t, err)
	golden.Assert(t, output)
}
