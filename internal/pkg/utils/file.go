package utils

import (
	"os"
)

// Mkdir ...
func Mkdir(dir string) bool {
	if dir == "" {
		return false
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return false
	}
	return true
}

// IsDir ...
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// IsExist returns whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
