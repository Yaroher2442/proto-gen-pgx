package helpers

import (
	"os"
	"path/filepath"
)

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDirIfNotExists(path string) error {
	exists, err := DirExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func FileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
