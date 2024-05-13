package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/mnogokotin/golang-packages/database/postgres"
	"github.com/mnogokotin/golang-packages/utils/e"
	"time"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/storage/tasks"
)

type Storage struct {
	*postgres.Postgres
}

func (s Storage) Delete(ctx context.Context, id int) error {
	result, err := s.Db.Exec(`delete from tasks 
	where id=$1`, id)
	rowsAffected, _ := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return e.Wrap("can't delete task by id", err)
	}

	return nil
}

func (s Storage) GetById(ctx context.Context, id int) (*domain.Task, error) {
	var task domain.Task
	err := s.Db.QueryRow(`select id, user_id, project_id, description, duration, date, created_at from tasks 
	where id=$1`, id).Scan(&task.ID, &task.UserId, &task.ProjectId, &task.Description, &task.Duration, &task.Date, &task.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Task{}, e.Wrap("can't get task by id", tasks.ErrTaskNotFound)
		}
		return &domain.Task{}, e.Wrap("can't get task by id", err)
	}

	return &task, nil
}

func (s Storage) GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error) {
	var tasks []*domain.Task

	rows, err := s.Db.Query(`select id, user_id, project_id, description, duration, date, created_at from tasks 
	where user_id=$1`, userId)
	if err != nil {
		return tasks, e.Wrap("can't get tasks list by userId", err)
	}

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.ID, &task.UserId, &task.ProjectId, &task.Description, &task.Duration, &task.Date, &task.CreatedAt)
		if err != nil {
			return tasks, e.Wrap("can't get tasks list by userId", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (s Storage) Store(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	var task_ domain.Task
	err := s.Db.QueryRow(`insert into tasks(user_id, project_id, description, duration, date, created_at)
	VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, task.UserId, task.ProjectId, task.Description, task.Duration, task.Date, time.Now()).Scan(&task_.ID)
	if err != nil {
		return nil, e.Wrap("can't store task", err)
	}

	return &task_, nil
}
