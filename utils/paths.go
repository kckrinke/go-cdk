package utils

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func FileWritable(path string) (writable bool) {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsRegular() {
			writable = unix.Access(path, unix.W_OK) == nil
		}
	}
	return
}

func IsFile(path string) (found bool) {
	if fi, err := os.Stat(path); err == nil {
		found = fi.Mode().IsRegular()
	}
	return
}

func IsDir(path string) (found bool) {
	if fi, err := os.Stat(path); err == nil {
		found = fi.Mode().IsDir()
	}
	return
}

func MakeDir(path string, perm os.FileMode) error {
	if IsDir(path) {
		return fmt.Errorf("directory exists")
	}
	if IsFile(path) {
		return fmt.Errorf("given path is a file")
	}
	return os.MkdirAll(path, perm)
}
