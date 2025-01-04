package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger initializes the zap logger
func InitLogger() {
	var err error
	
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
	zap.ReplaceGlobals(Logger)
}

// CloseLogger flushes any buffered log entries
func CloseLogger() {
	Logger.Sync()
}
