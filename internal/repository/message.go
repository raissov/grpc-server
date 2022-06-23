package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"server/internal/models"
	"strings"
	"time"
)

type UserRepository struct {
	db      *sql.DB
	timeout time.Duration
	logger  *zap.SugaredLogger
}

func newMessageRepo(db *sql.DB, timeout time.Duration, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:      db,
		timeout: timeout,
		logger:  logger,
	}
}

func (u *UserRepository) CreateMessages(messages []*models.Message) error {
	u.logger.Debug("CreateMessages - start")
	defer u.logger.Debug("CreateMessages - end")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ErrDBError
	}
	batch, err := tx.Prepare("INSERT INTO message (message)")
	if err != nil {
		return models.ErrDBError
	}

	for i := 0; i < len(messages); i++ {
		_, err = batch.Exec(
			messages[i].Message,
		)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return models.ErrDBError
			}
			return models.ErrDBError
		}
	}
	err = tx.Commit()
	if err != nil {
		return models.ErrDBError
	}
	return nil
}

func (u *UserRepository) CreateMessagesBatch(messages []*models.Message) error {
	u.logger.Debug("CreateMessages - start")
	defer u.logger.Debug("CreateMessages - end")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	valueStrings := make([]string, 0, len(messages))
	valueArgs := make([]interface{}, 0, len(messages))

	for index, message := range messages {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", index+1))
		valueArgs = append(valueArgs, message.Message)
	}
	tx, err := u.db.BeginTx(ctx, nil)
	query := fmt.Sprintf("INSERT INTO message (message) VALUES %s", strings.Join(valueStrings, ","))
	_, err = u.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return models.ErrDBError
	}
	if err := tx.Commit(); err != nil {
		return models.ErrDBError
	}
	return nil
}

func (u *UserRepository) GetMessages() ([]*models.MessageDB, error) {
	u.logger.Debug("GetMessages - start")
	defer u.logger.Debug("GetMessages - end")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	rows, err := u.db.QueryContext(ctx, "SELECT message, created_at FROM message LIMIT 5")
	if err != nil {
		return nil, models.ErrDBError
	}
	defer rows.Close()

	messages := make([]*models.MessageDB, 0)
	for rows.Next() {
		var message models.MessageDB
		if err := rows.Scan(&message.Message, &message.CreatedAt); err != nil {
			return nil, models.ErrDBError
		}
		messages = append(messages, &message)
	}
	if err := rows.Err(); err != nil {
		return nil, models.ErrDBError
	}
	return messages, nil
}
