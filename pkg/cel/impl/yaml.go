package impl

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"sigs.k8s.io/yaml"
)

// YamlImpl provides the concrete implementation for YAML operations
type YamlImpl struct{}

// Unmarshal parses YAML data and returns the result
func (YamlImpl) Unmarshal(data []byte) (any, error) {
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return result, nil
}

// UnmarshalFromFile reads and parses a YAML file from the local filesystem
func (YamlImpl) UnmarshalFromFile(path string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML from file %s: %w", path, err)
	}
	return result, nil
}

// UnmarshalFromURL fetches and parses a YAML file from a URL
func (YamlImpl) UnmarshalFromURL(url string) (any, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL %s: status code %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML from URL %s: %w", url, err)
	}
	return result, nil
}
