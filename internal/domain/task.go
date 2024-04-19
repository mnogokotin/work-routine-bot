package domain

import "time"

type Task struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id"`
	ProjectId   int       `json:"project_id"`
	Description string    `json:"description"`
	Duration    float64   `json:"duration"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}
