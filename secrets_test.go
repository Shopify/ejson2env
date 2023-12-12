package ejson2env

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestReadSecrets(t *testing.T) {
	var err error

	_, err = ReadAndExtractEnv("testdata/bad.ejson", "./key", TestKeyValue)
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
		include        []string
		expectedOutput string
	}{
		{
			name:           "ExportEnv",
			exportFunc:     ExportEnv,
			include:        make([]string, 0),
			expectedOutput: "export test_key='test value'\nexport test_key_2='test value 2'\nexport test_key_3='test value 3'\n",
		},
		{
			name:           "ExportQuiet",
			exportFunc:     ExportQuiet,
			include:        make([]string, 0),
			expectedOutput: "test_key='test value'\ntest_key_2='test value 2'\ntest_key_3='test value 3'\n",
		},
		{
			name:           "ExportInclude",
			exportFunc:     ExportEnv,
			include:        []string{"test_key", "test_key_3"},
			expectedOutput: "export test_key='test value'\nexport test_key_3='test value 3'\n",
		},
	}

	for _, test := range tests {
		err := ReadAndExportEnv("testdata/test-expected-usage.ejson", "./key", TestKeyValue, test.exportFunc, test.include)
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

	err = ReadAndExportEnv("bad.ejson", "./key", TestKeyValue, ExportEnv, make([]string, 0))
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}
