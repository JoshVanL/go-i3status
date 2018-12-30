package utils

import (
	"bytes"
	"io/ioutil"
)

func ReadFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(f), nil
}
