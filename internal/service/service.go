package service

import (
	"go.uber.org/zap"
	"server/config"
	"server/internal/models"
	"server/internal/repository"
)

type Service struct {
	VerifyToken
	Message
}

type Message interface {
	CreateMessages(messages []*models.Message) error
	CreateMessagesBatch(messages []*models.Message) error
	GetMessages() ([]*models.MessageDB, error)
}

type VerifyToken interface {
	VerifyToken(token string, url string) error
}

func NewService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *Service {
	return &Service{
		Message:     NewMessageService(repo.Message, cfg, logger),
		VerifyToken: NewVerifyTokenService(repo.Verify, cfg, logger),
	}
}
