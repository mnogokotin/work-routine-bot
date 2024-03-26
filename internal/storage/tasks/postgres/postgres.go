package postgres

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/utils/e"
	"time"
	"work-routine-bot/internal/domain"
)

type Storage struct {
	*postgres.Postgres
}

func (s Storage) Store(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	err := s.Db.QueryRow(`insert into tasks(user_id, project_id, description, duration, date, created_at)
	VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, task.UserId, task.ProjectId, task.Description, task.Duration, task.Date, time.Now()).Scan(task.ID)
	if err != nil {
		return nil, e.Wrap("can't store task", err)
	}

	return task, nil
}

func (s Storage) GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error) {
	var tasks []*domain.Task

	rows, err := s.Db.Query(`select id, user_id, project_id, description, duration, date, created_at from tasks 
	where user_id=$1`, userId)
	if err != nil {
		return tasks, e.Wrap("can't get work hours list by userId", err)
	}

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.ID, &task.UserId, &task.ProjectId, &task.Description, &task.Duration, &task.Date, &task.CreatedAt)
		if err != nil {
			return tasks, e.Wrap("can't get work hours list by userId", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}
