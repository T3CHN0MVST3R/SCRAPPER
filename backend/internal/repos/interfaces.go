package repos

import (
	"context"

	"scrapper/internal/dto"

	"github.com/google/uuid"
)

// ParserRepo представляет интерфейс для репозитория парсера
type ParserRepo interface {
	// CreateOperation создает новую операцию парсинга
	CreateOperation(ctx context.Context, url string) (uuid.UUID, error)

	// UpdateOperationStatus обновляет статус операции
	UpdateOperationStatus(ctx context.Context, operationID uuid.UUID, status dto.OperationStatus) error

	// GetOperationByID получает операцию по ID
	GetOperationByID(ctx context.Context, operationID uuid.UUID) (*dto.Operation, error)

	// SaveBlock сохраняет блок, найденный при парсинге
	SaveBlock(ctx context.Context, block *dto.Block) error

	// GetBlocksByOperationID получает все блоки по ID операции
	GetBlocksByOperationID(ctx context.Context, operationID uuid.UUID) ([]dto.Block, error)

	//GetAllTemplates получает все HTML теги для парсера блоков страницы
	GetAllTemplates(platform dto.Platform) ([]dto.BlockTemplate, error)
}
