package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := Environment{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// Skip directories and files with invalid names
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := dir + "/" + file.Name()
		fileContent, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer fileContent.Close()

		scanner := bufio.NewScanner(fileContent)
		// Pass for \n
		// Will false for len=0 files
		if scanner.Scan() {
			line := scanner.Text()

			// Replace null characters with newline and trim spaces and tabs from the end
			line = strings.ReplaceAll(line, string([]byte{0x00}), "\n")
			line = strings.TrimRight(line, " \t")

			env[file.Name()] = EnvValue{Value: line}
		} else {
			// If the file is empty, mark the variable for removal
			env[file.Name()] = EnvValue{NeedRemove: true}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return env, nil
}
