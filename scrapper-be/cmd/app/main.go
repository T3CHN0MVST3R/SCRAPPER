package main

import (
	"go.uber.org/fx"

	"go.uber.org/zap"

	"scrapper/api/v1"
	"scrapper/config"
	"scrapper/internal/app"
	"scrapper/internal/handlers"
	"scrapper/internal/repos"
	"scrapper/internal/services"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	application := fx.New(
		fx.Supply(logger),
		fx.Provide(
			config.NewConfig,
		),
		repos.Module,
		services.Module,
		handlers.Module,
		api.Module,
		app.Module,
	)

	application.Run()
}
