package api

import (
	"go.uber.org/fx"
)

// Module регистрирует зависимости для API
var Module = fx.Module("api",
	fx.Provide(
		NewRouter,
	),
)
