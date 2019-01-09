package main

import (
	"bytes"
	"errors"
	"testing"
)

var errBadRead = errors.New("broken reader")

type badReader struct {
}

func (b *badReader) Read(p []byte) (n int, err error) {
	return 0, errBadRead
}

func TestReadKey(t *testing.T) {
	buffer := bytes.NewBufferString(TestKeyValue)
	value, err := readKey(buffer)

	if nil != err {
		t.Errorf("shouldn't have returned an error, return: %s", err)
	}

	if value != TestKeyValue {
		t.Errorf("value does not match expected:\nvalue: \"%s\"\nexpected: \"%s\"", value, TestKeyValue)
	}
}

func TestReadBadKey(t *testing.T) {
	value, err := readKey(new(badReader))

	if errBadRead != err {
		t.Errorf("should have returned \"%s\" error, returned: %s", errBadRead, err)
	}

	if value != "" {
		t.Errorf("value should be empty.\nvalue: \"%s\"\n", value)
	}
}
