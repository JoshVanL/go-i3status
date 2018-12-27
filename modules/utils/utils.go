package utils

import (
	"bytes"
	"io/ioutil"
)

func ReadFile(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	return bytes.TrimSpace(f)
}
