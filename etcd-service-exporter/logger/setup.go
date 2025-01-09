package logger

import (
	"go.uber.org/zap"
)

// SetupGlobalLogger 用于初始化全局日志
func SetupGlobalLogger() error {
	logger, err := zap.NewProduction(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	)
	if err != nil {
		return err
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}(logger)

	zap.ReplaceGlobals(logger)
	return nil
}
