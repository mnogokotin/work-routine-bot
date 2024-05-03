package tasks

import (
	"context"
	"errors"
	"work-routine-bot/internal/domain"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*domain.Task, error)
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}
