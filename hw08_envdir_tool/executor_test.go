package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testEnvironment() Environment {
	return Environment{
		"BAR": EnvValue{"BAR", false},
	}
}

func TestRunCmd(t *testing.T) {
	t.Run("exit codes", func(t *testing.T) {
		code := RunCmd([]string{"/bin/bash", "testdata/exit.sh"}, testEnvironment())
		require.Equalf(t, 2, code, "Invalid code - %d", code)
	})

	t.Run("one argument", func(t *testing.T) {
		code := RunCmd([]string{"echo"}, testEnvironment())
		require.Equalf(t, 0, code, "Invalid code - %d", code)
	})
}
