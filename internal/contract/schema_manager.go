package contract

import (
	"go.uber.org/zap"
)

// SchemaManagerService is a file-based implementation of the Service interface
type SchemaManagerService struct {
	logger        *zap.Logger
	contractsPath string
}

// NewSchemaManagerService creates a new schema manager service instance
func NewSchemaManagerService(logger *zap.Logger, contractsPath string) Service {
	return &SchemaManagerService{
		logger:        logger,
		contractsPath: contractsPath,
	}
}
