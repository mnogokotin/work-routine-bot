package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/utils/e"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/storage/projects"
)

type Storage struct {
	*postgres.Postgres
}

func (s Storage) GetList(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project

	rows, err := s.Db.Query(`select id, name from projects 
	order by name`)
	if err != nil {
		return projects, e.Wrap("can't get projects list", err)
	}

	for rows.Next() {
		var project domain.Project
		err = rows.Scan(&project.ID, &project.Name)
		if err != nil {
			return projects, e.Wrap("can't get tasks list by userId", err)
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func (s Storage) GetByName(ctx context.Context, name string) (*domain.Project, error) {
	var project domain.Project
	err := s.Db.QueryRow(`select id, name FROM projects
	where name = $1`, name).Scan(&project.ID, &project.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Project{}, e.Wrap("can't get project by name", projects.ErrProjectNotFound)
		}
		return &domain.Project{}, e.Wrap("can't get project by name", err)
	}

	return &project, nil
}
