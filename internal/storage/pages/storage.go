package pages

import (
	"context"
	"errors"
	"work-routine-bot/internal/domain"
)

type Storage interface {
	Store(ctx context.Context, p *domain.Page) error
	GetRandom(ctx context.Context, username string) (*domain.Page, error)
	Remove(ctx context.Context, p *domain.Page) error
	IsExists(ctx context.Context, p *domain.Page) (bool, error)
}

var ErrNoStoredPages = errors.New("no stored pages")
