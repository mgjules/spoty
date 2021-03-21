package app

import "go.uber.org/fx"

var DefaultProviders = fx.Provide(
	ProvideConfig,
	ProvideCache,
	ProvideSpoty,
	ProvideServer,
)
