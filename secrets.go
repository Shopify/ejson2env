package ejson2env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Shopify/ejson"
)

// ExportFunction is implemented in exportSecrets as an easy way
// to select how secrets are exported
type ExportFunction func(io.Writer, map[string]string)

// output is a pointer to the io.Writer to use. This allows us to override
// stdout for testing purposes.
var output io.Writer = os.Stdout

// readSecrets reads the secrets for the passed filename and
// returns them as a map[string]interface{}.
func readSecrets(filename, keyDir, privateKey string) (map[string]interface{}, error) {
	secrets := make(map[string]interface{})

	decrypted, err := ejson.DecryptFile(filename, keyDir, privateKey)
	if nil != err {
		return secrets, err
	}

	decoder := json.NewDecoder(bytes.NewReader(decrypted))

	err = decoder.Decode(&secrets)
	return secrets, err
}

// IsEnvError returns true if the passed error is due to the environment
// being missing or not containing environment variables.
func IsEnvError(err error) bool {
	return (errNoEnv == err || errEnvNotMap == err)
}

// FilterEnv removes any key from the secrets that's not in include.
// If include is empty, return the secrets unchanged.
func FilterEnv(originalEnv map[string]string, include []string) (map[string]string, error) {
	if len(include) == 0 {
		return originalEnv, nil
	}

	filteredEnv := make(map[string]string, len(include))
	for _, key := range include {
		if value, exists := originalEnv[key]; exists {
			filteredEnv[key] = value
		} else {
			return map[string]string{}, fmt.Errorf("key not found in ejson file: %s", key)
		}
	}

	return filteredEnv, nil
}

// ReadAndExportEnv wraps the read, extract, and export steps. Returns
// an error if any step fails.
func ReadAndExportEnv(filename, keyDir, privateKey string, exportFunc ExportFunction, include []string) error {
	envValues, err := ReadAndExtractEnv(filename, keyDir, privateKey)

	if nil != err && !IsEnvError(err) {
		return fmt.Errorf("could not load environment from file: %s", err)
	}

	filteredEnv, err := FilterEnv(envValues, include)
	if nil != err {
		return err
	}

	exportFunc(output, filteredEnv)
	return nil
}
