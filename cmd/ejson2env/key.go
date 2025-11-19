package main

import (
	"io"
	"strings"
)

// readKey reads the contents of the passed reader, and
// strips any preceding or ending whitespace.
func readKey(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if nil != err {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}
