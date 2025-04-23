package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"scrapper/internal/dto"
	"scrapper/internal/services"
)

// parserHandler реализация ParserHandler
type parserHandler struct {
	logger  *zap.Logger
	service services.ParserService
}

// NewParserHandler создает новый экземпляр ParserHandler
func NewParserHandler(logger *zap.Logger, service services.ParserService) ParserHandler {
	return &parserHandler{
		logger:  logger,
		service: service,
	}
}

// ParseURL обрабатывает запрос на парсинг URL
func (h *parserHandler) ParseURL(w http.ResponseWriter, r *http.Request) {
	var req dto.ParseURLRequest

	// Декодируем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Проверяем URL
	if req.URL == "" {
		RespondWithError(w, http.StatusBadRequest, "URL is required")
		return
	}

	// Вызываем сервис для парсинга URL
	operationID, err := h.service.ParseURL(r.Context(), req.URL)
	if err != nil {
		h.logger.Error("Failed to parse URL", zap.Error(err))
		RespondWithError(w, http.StatusInternalServerError, "Failed to parse URL")
		return
	}

	// Формируем ответ
	response := dto.ParseURLResponse{
		OperationID: operationID,
	}

	RespondWithJSON(w, http.StatusOK, response)
}

// GetOperationResult обрабатывает запрос на получение результатов операции
func (h *parserHandler) GetOperationResult(w http.ResponseWriter, r *http.Request) {
	// Получаем ID операции из URL
	vars := mux.Vars(r)
	operationIDStr := vars["id"]

	// Проверяем ID операции
	operationID, err := uuid.Parse(operationIDStr)
	if err != nil {
		h.logger.Error("Invalid operation ID", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	// Вызываем сервис для получения результатов операции
	result, err := h.service.GetOperationResult(r.Context(), operationID)
	if err != nil {
		h.logger.Error("Failed to get operation result", zap.Error(err))
		RespondWithError(w, http.StatusInternalServerError, "Failed to get operation result")
		return
	}

	RespondWithJSON(w, http.StatusOK, result)
}

// ExportOperation обрабатывает запрос на экспорт результатов операции
func (h *parserHandler) ExportOperation(w http.ResponseWriter, r *http.Request) {
	// Получаем ID операции из URL
	vars := mux.Vars(r)
	operationIDStr := vars["id"]

	// Проверяем ID операции
	operationID, err := uuid.Parse(operationIDStr)
	if err != nil {
		h.logger.Error("Invalid operation ID", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	// Получаем формат экспорта из query параметров
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "excel" // По умолчанию Excel
	}

	// Проверяем формат
	if format != "excel" && format != "text" {
		RespondWithError(w, http.StatusBadRequest, "Invalid format. Supported formats: excel, text")
		return
	}

	// Вызываем сервис для экспорта операции
	content, filename, err := h.service.ExportOperation(r.Context(), operationID, format)
	if err != nil {
		h.logger.Error("Failed to export operation", zap.Error(err))
		RespondWithError(w, http.StatusInternalServerError, "Failed to export operation")
		return
	}

	// Устанавливаем заголовки для скачивания файла
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", getContentType(format))
	w.Header().Set("Content-Length", string(len(content)))

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// getContentType возвращает Content-Type в зависимости от формата
func getContentType(format string) string {
	switch format {
	case "excel":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "text":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

// RespondWithError отправляет клиенту ошибку в формате JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, dto.ErrorResponse{Error: message})
}

// RespondWithJSON отправляет клиенту данные в формате JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
