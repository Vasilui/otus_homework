package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrNoParams             = errors.New("not params in command line")
	ErrNoCommandForExecutor = errors.New("no command for executor")
	ErrReadDir              = errors.New("failed read dir")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrReadDir
	}

	for _, entry := range dirEntry {
		fileInfo, _ := entry.Info()
		if !validateEntry(fileInfo) {
			continue
		}

		file, err := os.Open(dir + "/" + entry.Name())
		if err != nil {
			fmt.Println("error open file ", err.Error())
			return nil, err
		}

		envVal, err := validatePotentialFile(fileInfo, file)
		if err != nil {
			continue
		}

		envs[fileInfo.Name()] = envVal
	}

	return envs, nil
}

func validateEntry(info os.FileInfo) bool {
	return !info.IsDir() && Executable(info.Mode()) && IsValidName(info.Name())
}

func Executable(mode os.FileMode) bool {
	return mode&0o111 == 0
}

func IsValidName(name string) bool {
	return !strings.Contains(name, "=")
}

func validatePotentialFile(info os.FileInfo, file *os.File) (EnvValue, error) {
	envVal := EnvValue{}

	if info.Size() == 0 {
		envVal.Value = ""
		envVal.NeedRemove = true
		return envVal, nil
	}

	if !info.Mode().IsRegular() {
		return envVal, errors.New("is not regular file")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res := bytes.ReplaceAll(scanner.Bytes(), []byte{0}, []byte{10})
		val := strings.TrimRight(string(res), "\t ")
		if val == "" {
			envVal.NeedRemove = true
			return envVal, nil
		}
		envVal.Value = val
		return envVal, nil
	}

	return envVal, nil
}

func ReadParams(params []string) (int, error) {
	if len(params) == 0 {
		return 1, ErrNoParams
	}

	if len(params) == 1 {
		return 1, ErrNoCommandForExecutor
	}

	env, err := ReadDir(params[0])
	if err != nil {
		return 1, err
	}

	code := RunCmd(params[1:], env)

	return code, nil
}
