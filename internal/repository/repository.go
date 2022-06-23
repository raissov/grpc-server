package repository

import (
	"database/sql"
	"go.uber.org/zap"
	"server/internal/models"
	"time"
)

type Repository struct {
	Message
	Verify
}

type Message interface {
	CreateMessages(messages []*models.Message) error
	CreateMessagesBatch(messages []*models.Message) error
	GetMessages() ([]*models.MessageDB, error)
}

type Verify interface {
	VerifyToken(token string, url string) error
}

func NewRepository(db *sql.DB, timeout time.Duration, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		Message: newMessageRepo(db, timeout, logger),
		Verify:  newVerifyTokenRepo(logger),
	}
}
