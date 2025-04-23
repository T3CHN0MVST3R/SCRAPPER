package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"scrapper/internal/services"
)

// crawlerHandler реализация CrawlerHandler
type crawlerHandler struct {
	logger  *zap.Logger
	service services.CrawlerService
}

// NewCrawlerHandler создает новый экземпляр CrawlerHandler
func NewCrawlerHandler(logger *zap.Logger, service services.CrawlerService) CrawlerHandler {
	return &crawlerHandler{
		logger:  logger,
		service: service,
	}
}

// CrawlURL обрабатывает запрос на обход URL и сбор ссылок
func (h *crawlerHandler) CrawlURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL      string `json:"url"`
		MaxDepth int    `json:"max_depth,omitempty"`
	}

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

	// Устанавливаем глубину обхода
	maxDepth := 2 // По умолчанию
	if req.MaxDepth > 0 {
		maxDepth = req.MaxDepth
	}

	// Проверяем, разрешен ли домен
	if !h.service.IsAllowedDomain(req.URL) {
		RespondWithError(w, http.StatusBadRequest, "Domain not allowed")
		return
	}

	// Обходим URL
	links, err := h.service.CrawlURL(r.Context(), req.URL, maxDepth)
	if err != nil {
		h.logger.Error("Failed to crawl URL", zap.Error(err))
		RespondWithError(w, http.StatusInternalServerError, "Failed to crawl URL")
		return
	}

	// Формируем ответ
	response := struct {
		URL   string   `json:"url"`
		Links []string `json:"links"`
		Count int      `json:"count"`
	}{
		URL:   req.URL,
		Links: links,
		Count: len(links),
	}

	RespondWithJSON(w, http.StatusOK, response)
}
