package task

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)

	CreateRecurrence(ctx context.Context, recurrence *taskdomain.TaskRecurrence) error
	GetRecurrenceByTaskID(ctx context.Context, taskID int64) (*taskdomain.TaskRecurrence, error)
	UpdateRecurrence(ctx context.Context, recurrence *taskdomain.TaskRecurrence) error
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type CreateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	RecurrenceType       taskdomain.RecurrenceType
	RecurrenceInterval   int
	RecurrenceDayOfMonth *int
	SpecificDates        string
	StartDate            time.Time
	EndDate              *time.Time
}

type UpdateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	HasRecurrence        bool
	RecurrenceType       taskdomain.RecurrenceType
	RecurrenceInterval   int
	RecurrenceDayOfMonth *int
	SpecificDates        string
	StartDate            time.Time
	EndDate              *time.Time
}
