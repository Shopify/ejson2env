package main

import (
	"bytes"
	"testing"
)

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
