package integration

import (
	"testing"

	"github.com/dcilke/golden"
	"github.com/stretchr/testify/require"
)

var filecases = []string{"strings", "ndjson", "pretty", "array", "mixed"}

func TestCLI(t *testing.T) {
	for _, tc := range filecases {
		t.Run(tc, func(t *testing.T) {
			output, err := hz(fn(tc), "--plain")
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Strict(t *testing.T) {
	for _, tc := range filecases {
		t.Run(tc, func(t *testing.T) {
			output, err := hz(fn(tc), "--plain", "--strict")
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Level(t *testing.T) {
	for _, tc := range filecases {
		for _, level := range []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"} {
			t.Run(tc+"/"+level, func(t *testing.T) {
				output, err := hz(fn(tc), "--plain", "--level", level)
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
			output, err := hz(tc, "--plain")
			require.NoError(t, err)
			golden.Assert(t, output)
		})
	}
}

func TestCLI_Flat(t *testing.T) {
	output, err := hz(fn("nested"), "--plain", "--flat")
	require.NoError(t, err)
	golden.Assert(t, output)
}

func TestCLI_Vert(t *testing.T) {
	output, err := hz(fn("nested"), "--plain", "--vertical")
	require.NoError(t, err)
	golden.Assert(t, output)
}

func TestCLI_Color(t *testing.T) {
	output, err := hz(fn("mixed"))
	require.NoError(t, err)
	golden.Assert(t, output)
}

func TestCLI_NoPin(t *testing.T) {
	output, err := hz(fn("mixed"), "--plain", "--no-pin")
	require.NoError(t, err)
	golden.Assert(t, output)
}
