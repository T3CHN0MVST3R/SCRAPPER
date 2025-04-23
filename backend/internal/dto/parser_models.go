package dto

import (
	"time"

	"github.com/google/uuid"
)

// OperationStatus представляет статус операции парсинга
type OperationStatus string

const (
	StatusPending    OperationStatus = "pending"
	StatusProcessing OperationStatus = "processing"
	StatusCompleted  OperationStatus = "completed"
	StatusError      OperationStatus = "error"
)

// BlockType представляет тип блока
type BlockType string

const (
	BlockTypeHeader BlockType = "header"
	BlockTypeFooter BlockType = "footer"
)

// Platform представляет платформу сайта
type Platform string

const (
	PlatformWordPress Platform = "wordpress"
	PlatformTilda     Platform = "tilda"
	PlatformBitrix    Platform = "bitrix"
	PlatformHTML5     Platform = "html5"
	PlatformUnknown   Platform = "unknown"
)

// Operation представляет операцию парсинга
type Operation struct {
	ID        uuid.UUID       `json:"id"`
	URL       string          `json:"url"`
	Status    OperationStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Block представляет блок, найденный при парсинге
type Block struct {
	ID          uuid.UUID   `json:"id"`
	OperationID uuid.UUID   `json:"operation_id"`
	BlockType   BlockType   `json:"block_type"`
	Platform    Platform    `json:"platform"`
	Content     interface{} `json:"content"`
	HTML        string      `json:"html"`
	CreatedAt   time.Time   `json:"created_at"`
}

// ParseURLRequest представляет запрос на парсинг URL
type ParseURLRequest struct {
	URL string `json:"url"`
}

// ParseURLResponse представляет ответ на запрос парсинга URL
type ParseURLResponse struct {
	OperationID uuid.UUID `json:"operation_id"`
}

// GetOperationResultRequest представляет запрос на получение результатов операции
type GetOperationResultRequest struct {
	OperationID uuid.UUID `json:"operation_id"`
}

// GetOperationResultResponse представляет ответ с результатами операции
type GetOperationResultResponse struct {
	Operation Operation `json:"operation"`
	Blocks    []Block   `json:"blocks"`
}

// ExportOperationRequest представляет запрос на экспорт результатов операции
type ExportOperationRequest struct {
	OperationID uuid.UUID `json:"operation_id"`
	Format      string    `json:"format"` // "excel" или "text"
}

// ExportOperationResponse представляет ответ на запрос экспорта
type ExportOperationResponse struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
