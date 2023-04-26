package golden

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/require"
)

const (
	EnvUpdate = "GOLDEN_UPDATE"
)

type Asserter func(t testing.TB, expected []byte, actual any)

func getPath(name string) string {
	return filepath.Join("testdata", "golden", name)
}

// AssertFunc fetches the test data from the golden directory,
// updating as necessary, and then calls the Asserter to compare
// the values
func AssertFunc(t testing.TB, name string, actual any, a Asserter) {
	t.Helper()
	file := getPath(name)
	actualb := ToBytes(t, actual)
	if os.Getenv(EnvUpdate) == "true" {
		Update(t, file, actualb)
	}

	expected, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		Update(t, file, actualb)
		expected = actualb
	} else {
		require.NoError(t, err)
	}

	a(t, expected, actual)
}

// Assert asserts the value for the given test
func Assert(t testing.TB, actual any) {
	t.Helper()
	AssertFunc(t, t.Name(), actual, func(t testing.TB, expected []byte, actual any) {
		t.Helper()
		a := ToBytes(t, actual)
		if !bytes.Equal(expected, a) {
			require.Fail(t, "Diff:\n"+diff.Diff(string(expected), string(a)))
		}
	})
}

// Subsert asserts on a sub assertion a test
func Subsert(t testing.TB, name string, actual any) {
	t.Helper()
	AssertFunc(t, filepath.Join(t.Name(), name), actual, func(t testing.TB, expected []byte, actual any) {
		t.Helper()
		a := ToBytes(t, actual)
		if !bytes.Equal(expected, a) {
			require.Fail(t, "Diff:\n"+diff.Diff(string(expected), string(a)))
		}
	})
}

// Update updates the golden file
func Update(t testing.TB, file string, data []byte) {
	t.Helper()
	t.Logf("updating golden file: %s", file)
	err := os.MkdirAll(path.Dir(file), 0755)
	require.NoError(t, err)
	err = os.WriteFile(file, data, 0644)
	require.NoError(t, err)
}

func ToBytes(t testing.TB, actual any) []byte {
	t.Helper()
	if actual == nil {
		return nil
	}
	switch v := actual.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	}
	a := litter.Sdump(actual)
	return []byte(a)
}
