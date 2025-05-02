package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// PostgresRepo реализация ParserRepo для PostgreSQL
type PostgresRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewParserRepo создает новый экземпляр ParserRepo
func NewParserRepo(db *sql.DB, logger *zap.Logger) ParserRepo {
	return &PostgresRepo{
		db:     db,
		logger: logger,
	}
}

// CreateOperation создает новую операцию парсинга
func (r *PostgresRepo) CreateOperation(ctx context.Context, url string) (uuid.UUID, error) {
	var operationID uuid.UUID

	query := `
	INSERT INTO operations (url, status)
	VALUES ($1, $2)
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, url, dto.StatusPending).Scan(&operationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create operation: %w", err)
	}

	return operationID, nil
}

// UpdateOperationStatus обновляет статус операции
func (r *PostgresRepo) UpdateOperationStatus(ctx context.Context, operationID uuid.UUID, status dto.OperationStatus) error {
	query := `
	UPDATE operations
	SET status = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, operationID)
	if err != nil {
		return fmt.Errorf("failed to update operation status: %w", err)
	}

	return nil
}

// GetOperationByID получает операцию по ID
func (r *PostgresRepo) GetOperationByID(ctx context.Context, operationID uuid.UUID) (*dto.Operation, error) {
	query := `
	SELECT id, url, status, created_at, updated_at
	FROM operations
	WHERE id = $1
	`

	var operation dto.Operation
	var status string

	err := r.db.QueryRowContext(ctx, query, operationID).Scan(
		&operation.ID,
		&operation.URL,
		&status,
		&operation.CreatedAt,
		&operation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("operation not found: %s", operationID)
		}
		return nil, fmt.Errorf("failed to get operation: %w", err)
	}

	operation.Status = dto.OperationStatus(status)

	return &operation, nil
}

// SaveBlock сохраняет блок, найденный при парсинге
func (r *PostgresRepo) SaveBlock(ctx context.Context, block *dto.Block) error {
	contentJSON, err := json.Marshal(block.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal block content: %w", err)
	}

	query := `
	INSERT INTO blocks (operation_id, block_type, platform, content, html)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		block.OperationID,
		block.BlockType,
		block.Platform,
		contentJSON,
		block.HTML,
	).Scan(&block.ID, &block.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}

	return nil
}

// GetBlocksByOperationID получает все блоки по ID операции
func (r *PostgresRepo) GetBlocksByOperationID(ctx context.Context, operationID uuid.UUID) ([]dto.Block, error) {
	query := `
	SELECT id, operation_id, block_type, platform, content, html, created_at
	FROM blocks
	WHERE operation_id = $1
	ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, operationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}
	defer rows.Close()

	var blocks []dto.Block

	for rows.Next() {
		var block dto.Block
		var blockType, platform string
		var contentJSON []byte

		err := rows.Scan(
			&block.ID,
			&block.OperationID,
			&blockType,
			&platform,
			&contentJSON,
			&block.HTML,
			&block.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan block: %w", err)
		}

		block.BlockType = dto.BlockType(blockType)
		block.Platform = dto.Platform(platform)

		var content map[string]interface{}
		if err := json.Unmarshal(contentJSON, &content); err != nil {
			return nil, fmt.Errorf("failed to unmarshal block content: %w", err)
		}

		block.Content = content
		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blocks: %w", err)
	}

	return blocks, nil
}

// GetAllTemplates получает все HTML теги для парсера блоков страницы
func (r *PostgresRepo) GetAllTemplates(platform dto.Platform) ([]dto.BlockTemplate, error) {
	var (
		rows *sql.Rows
		err  error
	)
	switch platform {
	case dto.PlatformWordPress:
		rows, err = r.db.Query("SELECT block_type, wordpress FROM block_templates ORDER BY (wordpress ->>'priority')::int")
	case dto.PlatformTilda:
		rows, err = r.db.Query("SELECT block_type, tilda FROM block_templates ORDER BY (tilda ->>'priority')::int")
	case dto.PlatformBitrix:
		rows, err = r.db.Query("SELECT block_type, bitrix FROM block_templates ORDER BY (bitrix ->>'priority')::int")
	case dto.PlatformHTML5:
		rows, err = r.db.Query("SELECT block_type, html5 FROM block_templates ORDER BY (html5 ->>'priority')::int")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var templates []dto.BlockTemplate

	for rows.Next() {
		var tmpl dto.BlockTemplate
		if err := rows.Scan(&tmpl.BlockType, &tmpl.HTMLTags); err != nil {
			return nil, fmt.Errorf("Error scanning block template: %w", err)
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}
