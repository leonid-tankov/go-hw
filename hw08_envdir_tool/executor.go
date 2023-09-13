package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	command := exec.Command(cmd[0], args...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return 1
			}
		}
		if value.Value == "" {
			continue
		}
		err := os.Setenv(key, value.Value)
		if err != nil {
			return 1
		}
	}
	command.Env = os.Environ()
	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
