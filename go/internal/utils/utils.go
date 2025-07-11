package utils

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func ExpandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return homeDir + path[1:], nil
	}
	return path, nil
}

func MarshalYAML(v any, header string) (string, error) {
	byteSlice, err := yaml.MarshalWithOptions(v, yaml.Indent(2))
	if err != nil {
		return "", fmt.Errorf("failed to marshal %T to YAML: %w", v, err)
	}
	if header != "" {
		return fmt.Sprintf("%s\n%s", header, string(byteSlice)), nil
	}
	return string(byteSlice), nil
}
