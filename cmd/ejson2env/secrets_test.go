package main

import (
	"strings"
	"testing"
)

func TestReadSecrets(t *testing.T) {
	var err error

	_, err = ReadSecrets("bad.ejson", "./key", TestKeyValue)
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}
