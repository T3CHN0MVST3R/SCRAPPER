package api

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"scrapper/internal/handlers"
)

// NewRouter создает и настраивает маршрутизатор
func NewRouter(
	logger *zap.Logger,
	parserHandler handlers.ParserHandler,
	downloaderHandler handlers.DownloaderHandler,
	crawlerHandler handlers.CrawlerHandler,
) *mux.Router {
	router := mux.NewRouter()

	// API маршруты
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Регистрируем маршруты парсера
	apiRouter.HandleFunc("/parse", parserHandler.ParseURL).Methods(http.MethodPost)
	apiRouter.HandleFunc("/operations/{id}", parserHandler.GetOperationResult).Methods(http.MethodGet)
	apiRouter.HandleFunc("/operations/{id}/export", parserHandler.ExportOperation).Methods(http.MethodGet)

	// Регистрируем маршруты загрузчика
	apiRouter.HandleFunc("/download/{id}", downloaderHandler.DownloadByID).Methods(http.MethodGet)
	apiRouter.HandleFunc("/formats", downloaderHandler.GetFormats).Methods(http.MethodGet)

	// Регистрируем маршруты краулера
	apiRouter.HandleFunc("/crawl", crawlerHandler.CrawlURL).Methods(http.MethodPost)

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
