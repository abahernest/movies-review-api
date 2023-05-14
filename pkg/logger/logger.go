package logger

import "go.uber.org/zap"

func InitLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	return logger, nil
}
