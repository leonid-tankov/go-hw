package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	environment := make(Environment, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if !fileInfo.Mode().IsRegular() {
			continue
		}
		if strings.Contains(fileInfo.Name(), "=") {
			continue
		}
		value, err := readFile(dir, fileInfo.Name())
		if err != nil {
			return nil, err
		}
		environment[fileInfo.Name()] = value
	}
	return environment, nil
}

func readFile(dir, file string) (EnvValue, error) {
	name := fmt.Sprintf("%s/%s", dir, file)
	content, err := os.ReadFile(name)
	if err != nil {
		return EnvValue{}, err
	}
	lines := strings.Split(string(content), "\n")
	content = bytes.ReplaceAll([]byte(lines[0]), []byte{0}, []byte("\n"))
	envValue := strings.TrimRight(string(content), " ")
	_, needRemove := os.LookupEnv(file)
	if len(envValue) == 0 {
		needRemove = true
	}
	return EnvValue{Value: envValue, NeedRemove: needRemove}, nil
}
