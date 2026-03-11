package logger

import "go.uber.org/zap"

func New(environment string) (*zap.Logger, error) {
	if environment == "prod" || environment == "production" {
		cfg := zap.NewProductionConfig()
		return cfg.Build()
	}

	cfg := zap.NewDevelopmentConfig()
	return cfg.Build()
}
