package files

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrNotExist = errors.New("file does not exist")
	ErrSymlink  = errors.New("file exists behind symlink and followLinks is false")
	ErrNotFile  = errors.New("path exists but is not a file")
)

func Exists(path string, followLinks bool) bool {
	var err error

	if followLinks {
		_, err = os.Stat(path)
	} else {
		_, err = os.Lstat(path)
	}

	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}

func IsFile(path string, followLinks bool) bool {
	var f os.FileInfo
	var err error

	if followLinks {
		f, err = os.Stat(path)
	} else {
		f, err = os.Lstat(path)
	}

	if err == nil {
		return !f.IsDir()
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}

func IsDir(path string, followLinks bool) bool {
	var f os.FileInfo
	var err error

	if followLinks {
		f, err = os.Stat(path)
	} else {
		f, err = os.Lstat(path)
	}

	if err == nil {
		return f.IsDir()
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}

func Read(path string, followLinks bool) ([]byte, error) {
	if !followLinks && IsFile(path, true) && !IsFile(path, false) {
		return []byte{}, ErrSymlink
	}

	if !Exists(path, true) {
		return []byte{}, ErrNotExist
	}

	if !IsFile(path, true) {
		return []byte{}, ErrNotFile
	}

	c, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("could not read file %s: %w", path, err)
	}

	return c, nil
}
