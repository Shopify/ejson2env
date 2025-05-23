package ejson2env

import (
	"bytes"
	"fmt"
	"testing"
)

const TestKeyValue = "2ed65dd6a16eab833cc4d2a860baa60042da34a58ac43855e8554ca87a5e557d"

func formatInvalid(received, expected string) string {
	return fmt.Sprintf("generated invalid export code: \n---\n%s\n---\nshould be: \n---\n%s\n---\n", received, expected)
}

func TestLoadSecrets(t *testing.T) {

	rawValues, err := readSecrets("testdata/test-expected-usage.ejson", "./key", TestKeyValue)
	if nil != err {
		t.Fatal(err)
	}

	envValues, err := ExtractEnv(rawValues)
	if nil != err {
		t.Fatal(err)
	}

	if "test value" != envValues["test_key"] {
		t.Error("Failed to decrypt")
	}
}

func TestLoadNoEnvSecrets(t *testing.T) {

	rawValues, err := readSecrets("testdata/test-public-key-only.ejson", "./key", TestKeyValue)
	if nil != err {
		t.Fatal(err)
	}

	_, err = ExtractEnv(rawValues)
	if errNoEnv != err {
		t.Fatal(err)
	}

	if !IsEnvError(err) {
		t.Fatalf("shouldn't have caused a failure: %s", err)
	}

}

func TestLoadBadEnvSecrets(t *testing.T) {

	rawValues, err := readSecrets("testdata/test-environment-string-not-object.ejson", "./key", TestKeyValue)
	if nil != err {
		t.Fatal(err)
	}

	_, err = ExtractEnv(rawValues)
	if errEnvNotMap != err {
		t.Fatal(err)
	}

	if !IsEnvError(err) {
		t.Fatalf("shouldn't have caused a failure: %s", err)
	}

}

func TestLoadUnderscoreEnvSecrets(t *testing.T) {

	rawValues, err := readSecrets("testdata/test-leading-underscore-env-key.ejson", "./key", TestKeyValue)
	if nil != err {
		t.Fatal(err)
	}

	envValues, err := ExtractEnv(rawValues)
	if nil != err {
		t.Fatal(err)
	}

	if "test value" != envValues["_test_key"] {
		t.Error("Failed to decrypt")
	}

}

func TestInvalidEnvironments(t *testing.T) {
	testGood := map[string]interface{}{
		"environment": map[string]interface{}{
			"test_key": "test_value",
		},
	}

	testBadNonMap := map[string]interface{}{
		"environment": "bad",
	}

	testBadInvalidKey := map[string]interface{}{
		"environment": map[string]interface{}{
			"invalid key": "test_value",
		},
	}

	var testNoEnv map[string]interface{}

	_, err := ExtractEnv(testBadNonMap)
	if nil == err {
		t.Errorf("no error when passed a non-map environment")
	} else if errEnvNotMap != err {
		t.Errorf("wrong error when passed a non-map environment: %s", err)
	}

	_, err = ExtractEnv(testBadInvalidKey)
	if nil == err {
		t.Errorf("no error when passed an environment with invalid key")
	} else if `invalid identifier as key in environment: "invalid key"` != err.Error() {
		t.Errorf("wrong error when passed an environment with invalid key: %s", err)
	}

	_, err = ExtractEnv(testNoEnv)
	if nil == err {
		t.Errorf("no error when passed a non-existent environment")
	} else if errNoEnv != err {
		t.Errorf("wrong error when passed a non-existent environment: %s", err)
	}

	_, err = ExtractEnv(testGood)
	if nil != err {
		t.Errorf("error when passed correctly formatted environment: %s", err)
	}

}

func TestEscaping(t *testing.T) {
	buf := bytes.Buffer{}

	testValues := map[string]string{
		"test": "test value'; echo dangerous; echo 'done",
	}

	ExportEnv(&buf, testValues)

	expectedOutput := "export test='test value'\"'\"'; echo dangerous; echo '\"'\"'done'\n"

	if expectedOutput != buf.String() {
		t.Fatal(formatInvalid(buf.String(), expectedOutput))
	}

}

func TestIdentifierPattern(t *testing.T) {
	key := "ALL_CAPS123"
	if !validIdentifierPattern.MatchString(key) {
		t.Errorf("key should match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "lowercase"
	if !validIdentifierPattern.MatchString(key) {
		t.Errorf("key should match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "a"
	if !validIdentifierPattern.MatchString(key) {
		t.Errorf("key should match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "_leading_underscore"
	if !validIdentifierPattern.MatchString(key) {
		t.Errorf("key should match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "1_leading_digit"
	if validIdentifierPattern.MatchString(key) {
		t.Errorf("key should not match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "contains whitespace"
	if validIdentifierPattern.MatchString(key) {
		t.Errorf("key should not match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "contains-dash"
	if validIdentifierPattern.MatchString(key) {
		t.Errorf("key should not match pattern %q: %q", validIdentifierPattern, key)
	}

	key = "contains_special_character;"
	if validIdentifierPattern.MatchString(key) {
		t.Errorf("key should not match pattern %q: %q", validIdentifierPattern, key)
	}
}
