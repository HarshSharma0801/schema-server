package contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

type SchemaService interface {
	Upload(ctx context.Context, schema *models.OpenAPI) error
	List(ctx context.Context) ([]string, error)
	Fetch(ctx context.Context, title string) (*models.OpenAPI, error)
	Download(ctx context.Context, title string) ([]byte, error)
}

type schemaService struct {
	logger   *zap.Logger
	basePath string
}

func NewSchemaService(logger *zap.Logger, basePath string) SchemaService {
	return &schemaService{logger: logger, basePath: basePath}
}

func (s *schemaService) schemaPath(title string) string {
	return filepath.Join(s.basePath, "schemas", fmt.Sprintf("%s.json", title))
}

func (s *schemaService) Upload(ctx context.Context, schema *models.OpenAPI) error {
	if schema.Info.Title == "" {
		return errors.New("schema title is required")
	}
	outputDir := filepath.Join(s.basePath, "schemas")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	filePath := s.schemaPath(schema.Info.Title)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return err
	}
	s.logger.Info("Uploaded schema", zap.String("title", schema.Info.Title), zap.String("path", filePath))
	return nil
}

func (s *schemaService) List(ctx context.Context) ([]string, error) {
	dir := filepath.Join(s.basePath, "schemas")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	titles := []string{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			titles = append(titles, file.Name()[:len(file.Name())-5])
		}
	}
	return titles, nil
}

func (s *schemaService) Fetch(ctx context.Context, title string) (*models.OpenAPI, error) {
	filePath := s.schemaPath(title)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var schema models.OpenAPI
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

func (s *schemaService) Download(ctx context.Context, title string) ([]byte, error) {
	filePath := s.schemaPath(title)
	return os.ReadFile(filePath)
}
