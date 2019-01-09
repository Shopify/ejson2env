package ejson2env

import (
	"errors"
	"fmt"
)

var errNoEnv = errors.New("environment is not set in ejson")
var errEnvNotMap = errors.New("environment is not a map[string]interface{}")

// ExtractEnv extracts the environment values from the map[string]interface{}
// containing all secrets, and returns a map[string]string containing the
// key value pairs. If there's an issue (the environment key doesn't exist, for
// example), returns an error.
func ExtractEnv(secrets map[string]interface{}) (map[string]string, error) {
	rawEnv, ok := secrets["environment"]
	if !ok {
		return nil, errNoEnv
	}

	envMap, ok := rawEnv.(map[string]interface{})
	if !ok {
		return nil, errEnvNotMap
	}

	envSecrets := make(map[string]string, len(envMap))

	for key, rawValue := range envMap {

		// Only export values that convert to strings properly.
		if value, ok := rawValue.(string); ok {
			envSecrets[key] = value
		}
	}

	return envSecrets, nil
}

// ReadAndExtractEnv wraps the read and extract steps. Returns
// a map[string]string containing the environment secrets to export.
func ReadAndExtractEnv(filename, keyDir, privateKey string) (map[string]string, error) {
	secrets, err := readSecrets(filename, keyDir, privateKey)
	if nil != err {
		return map[string]string{}, fmt.Errorf("could not load ejson file: %s", err)
	}

	return ExtractEnv(secrets)
}
