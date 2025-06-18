package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"schema-server/internal/models"
)

// File I/O operations

// ensureDirectories creates the necessary directory structure for contracts
func (s *SchemaManagerService) ensureDirectories() error {
	dirs := []string{
		filepath.Join(s.contractsPath, "mocks"),
		filepath.Join(s.contractsPath, "tests"),
		filepath.Join(s.contractsPath, "generated"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// readOpenAPIFile reads and parses an OpenAPI file
func (s *SchemaManagerService) readOpenAPIFile(path string) (*models.OpenAPI, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var openAPI models.OpenAPI
	if err := json.Unmarshal(data, &openAPI); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI: %w", err)
	}

	return &openAPI, nil
}

// writeOpenAPIFile writes an OpenAPI spec to a file
func (s *SchemaManagerService) writeOpenAPIFile(path string, openAPI *models.OpenAPI) error {
	data, err := json.MarshalIndent(openAPI, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// downloadExternalSchemas simulates downloading schemas from external sources
func (s *SchemaManagerService) downloadExternalSchemas() error {
	s.logger.Info("Simulating external schema download")
	// This would typically involve HTTP calls to external APIs
	return nil
}

// copyDirectory recursively copies a directory
func (s *SchemaManagerService) copyDirectory(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := s.copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := s.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func (s *SchemaManagerService) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
