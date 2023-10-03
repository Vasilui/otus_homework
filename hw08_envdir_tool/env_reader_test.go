package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		path, _ := os.Getwd()
		pathToDir := filepath.Join(path, "/testdata/invalid")

		envMap, errReadDir := ReadDir(pathToDir)

		require.Empty(t, envMap)
		require.NoError(t, errReadDir)
	})

	t.Run("test success", func(t *testing.T) {
		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir("testdata/env")

		fmt.Println(env)

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})

	t.Run("test  not use invalid name env", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		fmt.Println(env)

		require.NoError(t, err)
		_, ok := env["NOTUSE="]
		require.False(t, ok)
	})
}
