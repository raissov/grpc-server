package service

import (
	"go.uber.org/zap"
	"server/config"
	"server/internal/repository"
)

type VerifyTokenService struct {
	verifyToken repository.Verify
	cfg         *config.Configs
	logger      *zap.SugaredLogger
}

func (v *VerifyTokenService) VerifyToken(token string, url string) error {
	return v.verifyToken.VerifyToken(token, url)
}

func NewVerifyTokenService(verifyToken repository.Verify, cfg *config.Configs, logger *zap.SugaredLogger) *VerifyTokenService {
	return &VerifyTokenService{
		verifyToken: verifyToken,
		cfg:         cfg,
		logger:      logger,
	}
}
