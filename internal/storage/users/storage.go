package users

import (
	"context"
	"errors"
	"work-routine-bot/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Store(ctx context.Context, user *domain.User) (*domain.User, error)
}
