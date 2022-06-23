package service

import (
	"go.uber.org/zap"
	"server/config"
	"server/internal/models"
	"server/internal/repository"
)

type MessageService struct {
	messageRepo repository.Message
	cfg         *config.Configs
	logger      *zap.SugaredLogger
}

func (m MessageService) GetMessages() ([]*models.MessageDB, error) {
	return m.messageRepo.GetMessages()
}

func (m MessageService) CreateMessagesBatch(messages []*models.Message) error {
	return m.messageRepo.CreateMessagesBatch(messages)
}

func (m MessageService) CreateMessages(messages []*models.Message) error {
	return m.messageRepo.CreateMessages(messages)
}

func NewMessageService(messageRepo repository.Message, cfg *config.Configs, logger *zap.SugaredLogger) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		cfg:         cfg,
		logger:      logger,
	}
}
