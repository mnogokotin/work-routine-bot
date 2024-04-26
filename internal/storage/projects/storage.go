package projects

import (
	"context"
	"errors"
	"work-routine-bot/internal/domain"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type Storage interface {
	GetByName(ctx context.Context, name string) (*domain.Project, error)
	GetList(ctx context.Context) ([]*domain.Project, error)
}
