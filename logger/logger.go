package logger

import (
	"fmt"

	"github.com/JulesMike/spoty/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module exported for initialising a new Logger.
var Module = fx.Options(
	fx.Provide(New),
)

// Logger is a simple wrapper around zap.SugaredLogger.
type Logger struct {
	*otelzap.SugaredLogger
}

// New creates a new Logger.
func New(cfg *config.Config) (*Logger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	if cfg.Prod {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	otellogger := otelzap.New(logger)

	return &Logger{otellogger.Sugar()}, nil
}
