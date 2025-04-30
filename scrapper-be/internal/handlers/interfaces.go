package handlers

import (
	"net/http"
)

// ParserHandler представляет интерфейс для обработчика парсера
type ParserHandler interface {
	// ParseURL обрабатывает запрос на парсинг URL
	ParseURL(w http.ResponseWriter, r *http.Request)

	// GetOperationResult обрабатывает запрос на получение результатов операции
	GetOperationResult(w http.ResponseWriter, r *http.Request)

	// ExportOperation обрабатывает запрос на экспорт результатов операции
	ExportOperation(w http.ResponseWriter, r *http.Request)
}

// DownloaderHandler представляет интерфейс для обработчика загрузчика
type DownloaderHandler interface {
	// DownloadByID обрабатывает запрос на загрузку файлов по ID операции
	DownloadByID(w http.ResponseWriter, r *http.Request)

	// GetFormats обрабатывает запрос на получение доступных форматов
	GetFormats(w http.ResponseWriter, r *http.Request)
}

// CrawlerHandler представляет интерфейс для обработчика краулера
type CrawlerHandler interface {
	// CrawlURL обрабатывает запрос на обход URL и сбор ссылок
	CrawlURL(w http.ResponseWriter, r *http.Request)
}
