package handlers

import (
	"go.uber.org/fx"
)

// Module регистрирует зависимости для обработчиков
var Module = fx.Module("handlers",
	fx.Provide(
		NewParserHandler,
		NewDownloaderHandler,
		NewCrawlerHandler,
	),
)
