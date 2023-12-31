// internal/util/util.go

package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// GetEnv receive a slice of strings and return a map of strings of env variable and its value and an error
func GetEnv(envs []string) (map[string]string, error) {
	envMap := make(map[string]string)
	for _, env := range envs {
		envValue, ok := os.LookupEnv(env)
		if !ok {
			return nil, fmt.Errorf("env variable %s not found", env)
		}
		envMap[env] = envValue
	}
	return envMap, nil
}

// ParseBody receive a io.ReadCloser and a interface and return an error
func ParseBody(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}
