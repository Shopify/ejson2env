package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Shopify/ejson"
)

// ReadSecrets reads the secrets for the passed filename and
// returns them as a map[string]interface{}.
func ReadSecrets(filename, keyDir, privateKey string) (map[string]interface{}, error) {
	secrets := make(map[string]interface{})

	decrypted, err := ejson.DecryptFile(filename, keyDir, privateKey)
	if nil != err {
		return secrets, err
	}

	decoder := json.NewDecoder(bytes.NewReader(decrypted))

	err = decoder.Decode(&secrets)
	if nil != err {
		return secrets, err
	}

	return secrets, nil
}

// isFailure returns true if the passed error should prompt a
// failure.
func isFailure(err error) bool {
	return (nil != err && errNoEnv != err && errEnvNotMap != err)
}

// exportSecrets wraps the read, extract, and export steps. Returns
// an error if any step fails.
func exportSecrets(filename, keyDir, privateKey string) error {
	secrets, err := ReadSecrets(filename, keyDir, privateKey)
	if nil != err {
		return fmt.Errorf("could not load ejson file: %s", err)
	}

	envValues, err := ExtractEnv(secrets)
	if isFailure(err) {
		return fmt.Errorf("could not load environment from file: %s", err)
	}

	ExportEnv(os.Stdout, envValues)
	return nil
}
