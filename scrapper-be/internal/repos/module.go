package repos

import (
	"database/sql"
	"scrapper/config"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module регистрирует зависимости для репозиториев
var Module = fx.Module("repos",
	fx.Provide(
		NewPostgresConnection,
		NewParserRepo,
	),
)

// NewPostgresConnection создает подключение к базе данных PostgreSQL
func NewPostgresConnection(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Database.GetPostgresDSN())
	if err != nil {
		return nil, err
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Пингуем базу для проверки соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to PostgreSQL database")
	return db, nil
}
