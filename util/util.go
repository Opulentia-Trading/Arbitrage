package util

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"
)

// Returns the filename of the current file. Equivalent to __filename in other languages
func CurFilename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}

	return filename, nil
}

// Returns the directory of the current file. Equivalent to __dirname in other languages
func CurDirname() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}

	return filepath.Dir(filename), nil
}

// Pretty prints an object
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
