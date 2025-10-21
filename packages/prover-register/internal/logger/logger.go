package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with the specified configuration
func NewLogger(json bool, debug bool) (*zap.SugaredLogger, error) {
	var config zap.Config

	if json {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Development = true
	}

	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

// NewJSONLogger creates a production JSON logger
func NewJSONLogger() *zap.SugaredLogger {
	logger, _ := NewLogger(true, false)
	return logger
}

// NewDevelopmentLogger creates a development console logger
func NewDevelopmentLogger() *zap.SugaredLogger {
	logger, _ := NewLogger(false, true)
	return logger
}
