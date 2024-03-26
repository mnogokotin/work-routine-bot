package tasks

import (
	"context"
	"work-routine-bot/internal/domain"
)

type Storage interface {
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}
