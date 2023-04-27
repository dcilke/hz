package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const ()

var command string

func TestMain(m *testing.M) {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println(fmt.Errorf("getting go path: %w", err))
		os.Exit(1)
	}
	command = filepath.Join(root, "..", "bin", "test", "hz")
	os.Exit(m.Run())
}

func hz(args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=../.covdata")
	return cmd.CombinedOutput()
}

func fn(file string) string {
	return filepath.Join("testdata", "samples", file)
}
