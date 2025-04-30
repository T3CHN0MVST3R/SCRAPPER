package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"scrapper/config"
)

// Module регистрирует зависимости для приложения
var Module = fx.Module("app",
	fx.Invoke(
		StartServer,
	),
)

// StartServer запускает HTTP-сервер
func StartServer(lifecycle fx.Lifecycle, router *mux.Router, cfg *config.Config, logger *zap.Logger) {
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server", zap.String("port", cfg.Server.Port))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()

			// Обработка сигналов для грациозного завершения
			go func() {
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
				sig := <-sigChan
				logger.Info("Received signal", zap.String("signal", sig.String()))

				// Создаем контекст с таймаутом для завершения
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := server.Shutdown(shutdownCtx); err != nil {
					logger.Error("Server shutdown error", zap.Error(err))
				}

				logger.Info("Server gracefully stopped")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
