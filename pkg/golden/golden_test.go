package golden_test

import (
	"testing"

	"github.com/dcilke/hz/pkg/golden"
	"github.com/stretchr/testify/require"
)

func TestToBytes(t *testing.T) {
	testcases := map[string]struct {
		in  any
		out []byte
	}{
		"string": {"hello", []byte("hello")},
		"number": {42, []byte("42")},
		"bytes":  {[]byte{1, 2, 3}, []byte{1, 2, 3}},
		"object": {struct{ A string }{"hello"}, []byte("struct { A string }{\n  A: \"hello\",\n}")},
		"array":  {[]string{"hello", "world"}, []byte("[]string{\n  \"hello\",\n  \"world\",\n}")},
		"nil":    {nil, nil},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			out := golden.ToBytes(t, tc.in)
			require.Equal(t, tc.out, out)
		})
	}
}

func TestAssert(t *testing.T) {
	testcases := map[string]struct {
		actual any
	}{
		"string": {"hello"},
		"number": {42},
		"bytes":  {[]byte{1, 2, 3}},
		"object": {struct{ A string }{"hello"}},
		"array":  {[]string{"hello", "world"}},
		"nil":    {nil},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			golden.Assert(t, tc.actual)
		})
	}
}

func TestSubsert(t *testing.T) {
	testcases := map[string]struct {
		actual any
	}{
		"string": {"hello"},
		"number": {42},
		"bytes":  {[]byte{1, 2, 3}},
		"object": {struct{ A string }{"hello"}},
		"array":  {[]string{"hello", "world"}},
		"nil":    {nil},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			golden.Subsert(t, "1", tc.actual)
			golden.Subsert(t, "2", tc.actual)
			golden.Subsert(t, "3", tc.actual)
		})
	}
}
