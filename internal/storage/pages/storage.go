package pages

import (
	"context"
	"errors"
	"work-routine-bot/internal/domain"
)

type Storage interface {
	Save(ctx context.Context, p *domain.Page) error
	PickRandom(ctx context.Context, userName string) (*domain.Page, error)
	Remove(ctx context.Context, p *domain.Page) error
	IsExists(ctx context.Context, p *domain.Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")
