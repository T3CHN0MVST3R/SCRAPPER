package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"scrapper/internal/services"
	"bytes"
	"fmt"
	"time"
	"io"
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
	// Получаем ID операции из URL
	vars := mux.Vars(r)
	operationIDStr := vars["id"]

	operationID, err := uuid.Parse(operationIDStr)
	if err != nil {
		h.logger.Error("Invalid operation ID", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	// Получаем формат из query параметров
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "excel" // По умолчанию
	}

	// Проверяем поддерживаемые форматы
	formats := h.service.GetAvailableFormats()
	formatSupported := false
	for _, f := range formats {
		if f == format {
			formatSupported = true
			break
		}
	}

	if !formatSupported {
		RespondWithError(w, http.StatusBadRequest, "Unsupported format")
		return
	}

	// Получаем данные файла
	data, filename, err := h.service.DownloadByOperationID(r.Context(), operationID, format)
	if err != nil {
		h.logger.Error("Failed to generate file", zap.Error(err))
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate file")
		return
	}

	// Используем буферизированную отправку
	buffer := bytes.NewBuffer(data)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", getContentType(format))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	// Используем io.Copy для потоковой передачи
	if _, err := io.Copy(w, buffer); err != nil {
		h.logger.Error("Failed to send file", zap.Error(err))
		return
	}
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

func getContentType(format string) string {
	switch format {
	case "excel":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "pdf":
		return "application/pdf"
	case "text":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
