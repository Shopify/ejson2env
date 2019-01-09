package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestReadSecrets(t *testing.T) {
	var err error

	_, err = ReadSecrets("../../testdata/bad.ejson", "./key", TestKeyValue)
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}

func TestReadAndExportEnv(t *testing.T) {
	outputBuffer := new(bytes.Buffer)
	output = outputBuffer

	// ensure that output returns to os.Stdout
	defer func() {
		output = os.Stdout
	}()

	tests := []struct {
		name           string
		exportFunc     ExportFunction
		expectedOutput string
	}{
		{
			name:           "ExportEnv",
			exportFunc:     ExportEnv,
			expectedOutput: "export test_key='test value'\n",
		},
		{
			name:           "ExportQuiet",
			exportFunc:     ExportQuiet,
			expectedOutput: "test_key='test value'\n",
		},
	}

	for _, test := range tests {
		err := exportSecrets("../../testdata/test.ejson", "./key", TestKeyValue, test.exportFunc)
		if nil != err {
			t.Errorf("testing %s failed: %s", test.name, err)
			continue
		}

		actualOutput := outputBuffer.String()

		if test.expectedOutput != actualOutput {
			t.Error(formatInvalid(actualOutput, test.expectedOutput))
		}
		outputBuffer.Reset()
	}
}

func TestReadAndExportEnvWithBadEjson(t *testing.T) {
	var err error

	outputBuffer := new(bytes.Buffer)
	output = outputBuffer

	// ensure that output returns to os.Stdout
	defer func() {
		output = os.Stdout
	}()

	err = exportSecrets("../../testdata/bad.ejson", "./key", TestKeyValue, ExportEnv)
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}
