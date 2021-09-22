package helper

import (
	"io"
	"os"
)

func Copy(from, to string) error {
	toFile, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0o777)
	if err != nil {
		return err
	}
	defer toFile.Close()

	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = io.Copy(toFile, fromFile)
	return err
}
