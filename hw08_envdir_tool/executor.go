package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(command []string, env Environment) (returnCode int) {
	if command == nil || env == nil {
		return 1
	}
	cmd := exec.Command(command[0], command[1:]...) //nolint:gosec
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				continue
			}
		}

		err := os.Setenv(key, val.Value)
		if err != nil {
			return 1
		}
	}

	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return cmd.ProcessState.ExitCode()
}
