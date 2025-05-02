package services

import (
	"context"

	"scrapper/internal/dto"

	"github.com/google/uuid"
)

// PlatformService представляет интерфейс для сервиса конкретной платформы
type PlatformService interface {
	// DetectPlatform проверяет, соответствует ли страница данной платформе
	DetectPlatform(html string) bool

	// ParseHeader парсит шапку сайта
	ParseHeader(html string) (*dto.Block, error)

	// ParseFooter парсит подвал сайта
	ParseFooter(html string) (*dto.Block, error)
}

// ParserService представляет интерфейс для сервиса парсинга
type ParserService interface {
	// ParseURL парсит URL и сохраняет результаты в базу данных
	ParseURL(ctx context.Context, url string) (uuid.UUID, error)

	// GetOperationResult получает результаты операции по ID
	GetOperationResult(ctx context.Context, operationID uuid.UUID) (*dto.GetOperationResultResponse, error)

	// ExportOperation экспортирует результаты операции в файл
	ExportOperation(ctx context.Context, operationID uuid.UUID, format string) ([]byte, string, error)

	// DetectPlatform определяет платформу сайта по HTML
	DetectPlatform(html string) dto.Platform
}

// WordPressService представляет интерфейс для сервиса WordPress
type WordPressService interface {
	PlatformService
}

// TildaService представляет интерфейс для сервиса Tilda
type TildaService interface {
	PlatformService
}

// BitrixService представляет интерфейс для сервиса Bitrix
type BitrixService interface {
	PlatformService
}

// HTML5Service представляет интерфейс для сервиса HTML5
type HTML5Service interface {
	PlatformService
	ParseAndClassifyPage(html string, templates []dto.BlockTemplate) ([]*dto.Block, error)
}

// DownloaderService представляет интерфейс для сервиса загрузки файлов
type DownloaderService interface {
	// DownloadByOperationID загружает файлы по ID операции и сохраняет их в указанный путь
	DownloadByOperationID(ctx context.Context, operationID uuid.UUID, path string) error

	// GetAvailableFormats возвращает список доступных форматов для загрузки
	GetAvailableFormats() []string

	// DownloadByOperationIDWithFormat загружает файлы в указанном формате
	DownloadByOperationIDWithFormat(ctx context.Context, operationID uuid.UUID, format string, path string) error
}

type CrawlerService interface {
	// CrawlURL обходит URL и собирает ссылки
	CrawlURL(ctx context.Context, url string, maxDepth int) ([]string, error)

	// IsAllowedDomain проверяет, разрешен ли домен для обхода
	IsAllowedDomain(url string) bool

	// SetUserAgent устанавливает User-Agent для запросов
	SetUserAgent(userAgent string)

	// SetMaxDepth устанавливает максимальную глубину обхода
	SetMaxDepth(depth int)
}
