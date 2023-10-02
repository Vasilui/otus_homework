package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "foo"}
		env["BAR"] = EnvValue{Value: "bar"}
		env["HELLO"] = EnvValue{Value: "\"hello\""}

		code := RunCmd(cmd, env)

		require.Equal(t, 0, code)
	})

	t.Run("test invalid cmd", func(t *testing.T) {
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "foo"}
		env["BAR"] = EnvValue{Value: "bar"}

		code := RunCmd(nil, env)

		require.Equal(t, 1, code)
	})

	t.Run("test invalid env", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		code := RunCmd(cmd, nil)

		require.Equal(t, 1, code)
	})
}
