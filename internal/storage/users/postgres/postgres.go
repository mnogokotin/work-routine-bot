package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/utils/e"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/storage/users"
)

type Storage struct {
	*postgres.Postgres
}

func (s Storage) Store(ctx context.Context, user *domain.User) (*domain.User, error) {
	var user_ domain.User
	err := s.Db.QueryRow(`insert into users(username)
	VALUES($1) RETURNING id, username`, user.Username).Scan(&user_.ID, &user_.Username)
	if err != nil {
		return nil, e.Wrap("can't store task", err)
	}

	return &user_, nil
}

func (s Storage) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := s.Db.QueryRow(`select id, username FROM users
	where username = $1`, username).Scan(&user.ID, &user.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, e.Wrap("can't get user by username", users.ErrUserNotFound)
		}
		return &domain.User{}, e.Wrap("can't get user by username", err)
	}

	return &user, nil
}
