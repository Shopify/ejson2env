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

// ReadAndExportEnv wraps the read, extract, and export steps. Returns
// an error if any step fails.
func ReadAndExportEnv(filename, keyDir, privateKey string, exportFunc ExportFunction) error {
	envValues, err := ReadAndExtractEnv(filename, keyDir, privateKey)

	if nil != err && !IsEnvError(err) {
		return fmt.Errorf("could not load environment from file: %s", err)
	}

	exportFunc(output, envValues)
	return nil
}
