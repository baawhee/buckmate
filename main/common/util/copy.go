package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func RemoveAllFromDirectory(dir string) error {
	return os.RemoveAll(dir + "/*")
}
func CopyDirectory(scrDir string, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	if err := CreateIfNotExists(dest, 0755); err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := out.Close(); cerr != nil {
			err = cerr
		}
	}()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := in.Close(); cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
