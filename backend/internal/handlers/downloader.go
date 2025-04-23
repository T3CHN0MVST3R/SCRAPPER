package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"scrapper/internal/services"
)

// downloaderHandler реализация DownloaderHandler
type downloaderHandler struct {
	logger  *zap.Logger
	service services.DownloaderService
}

// NewDownloaderHandler создает новый экземпляр DownloaderHandler
func NewDownloaderHandler(logger *zap.Logger, service services.DownloaderService) DownloaderHandler {
	return &downloaderHandler{
		logger:  logger,
		service: service,
	}
}

// DownloadByID обрабатывает запрос на загрузку файлов по ID операции
func (h *downloaderHandler) DownloadByID(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// GetFormats обрабатывает запрос на получение доступных форматов
func (h *downloaderHandler) GetFormats(w http.ResponseWriter, r *http.Request) {
	formats := h.service.GetAvailableFormats()

	response := struct {
		Formats []string `json:"formats"`
	}{
		Formats: formats,
	}

	RespondWithJSON(w, http.StatusOK, response)
}
