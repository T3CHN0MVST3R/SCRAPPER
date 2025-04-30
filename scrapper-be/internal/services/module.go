package services

import (
	"go.uber.org/fx"
	"scrapper/config"
)

// Module регистрирует зависимости для сервисов
var Module = fx.Module("services",
	fx.Provide(
		NewWordPressService,
		NewTildaService,
		NewBitrixService,
		NewHTML5Service,
		NewParserService,
		NewDownloaderService,
		func(cfg *config.Config) []string {
			return cfg.Scraper.AllowedDomains
		},
		NewCrawlerService,
	),
)
