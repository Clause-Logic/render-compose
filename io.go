package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WriteToFile writes the blueprint to a YAML file
func (bp *Blueprint) WriteToFile(path string) error {
	if bp == nil {
		return fmt.Errorf("blueprint is nil")
	}

	// Validate blueprint before writing
	if errors := ValidateBlueprint(bp); len(errors) > 0 {
		return fmt.Errorf("blueprint validation failed: %s", strings.Join(errors, "; "))
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Serialize to YAML
	data, err := yaml.Marshal(bp)
	if err != nil {
		return fmt.Errorf("failed to marshal blueprint to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// WriteRenderYAML writes the blueprint to render.yaml in the current directory
func (bp *Blueprint) WriteRenderYAML() error {
	return bp.WriteToFile("render.yaml")
}

// WriteRenderYAMLTo writes the blueprint to render.yaml in the specified directory
func (bp *Blueprint) WriteRenderYAMLTo(dir string) error {
	return bp.WriteToFile(filepath.Join(dir, "render.yaml"))
}

// ToYAMLString converts the blueprint to a YAML string
func (bp *Blueprint) ToYAMLString() (string, error) {
	if bp == nil {
		return "", fmt.Errorf("blueprint is nil")
	}

	data, err := yaml.Marshal(bp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal blueprint to YAML: %w", err)
	}

	return string(data), nil
}

// ToYAMLBytes converts the blueprint to YAML bytes
func (bp *Blueprint) ToYAMLBytes() ([]byte, error) {
	if bp == nil {
		return nil, fmt.Errorf("blueprint is nil")
	}

	data, err := yaml.Marshal(bp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blueprint to YAML: %w", err)
	}

	return data, nil
}

// LoadFromFile loads a blueprint from a YAML file
func LoadFromFile(path string) (*Blueprint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var bp Blueprint
	if err := yaml.Unmarshal(data, &bp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML from %s: %w", path, err)
	}

	return &bp, nil
}

// LoadRenderYAML loads a blueprint from render.yaml in the current directory
func LoadRenderYAML() (*Blueprint, error) {
	return LoadFromFile("render.yaml")
}

// LoadRenderYAMLFrom loads a blueprint from render.yaml in the specified directory
func LoadRenderYAMLFrom(dir string) (*Blueprint, error) {
	return LoadFromFile(filepath.Join(dir, "render.yaml"))
}

// WriteWithBackup writes the blueprint to a file, creating a backup if the file exists
func (bp *Blueprint) WriteWithBackup(path string) error {
	// Create backup if file exists
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".backup"
		if err := copyFile(path, backupPath); err != nil {
			return fmt.Errorf("failed to create backup %s: %w", backupPath, err)
		}
	}

	return bp.WriteToFile(path)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}