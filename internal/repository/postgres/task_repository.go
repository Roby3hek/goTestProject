package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return scanTask(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	return scanTask(row)
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.ID,
	)

	return scanTask(row)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (r *Repository) CreateRecurrence(ctx context.Context, recurrence *taskdomain.TaskRecurrence) error {
	const query = `
		INSERT INTO task_recurrences (
			task_id, recurrence_type, recurrence_interval, recurrence_day_of_month,
			specific_dates, start_date, end_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		recurrence.TaskID,
		recurrence.RecurrenceType,
		recurrence.RecurrenceInterval,
		recurrence.RecurrenceDayOfMonth,
		recurrence.SpecificDates,
		recurrence.StartDate,
		recurrence.EndDate,
		recurrence.CreatedAt,
		recurrence.UpdatedAt,
	)
	return err
}

func (r *Repository) GetRecurrenceByTaskID(ctx context.Context, taskID int64) (*taskdomain.TaskRecurrence, error) {
	const query = `
		SELECT task_id, recurrence_type, recurrence_interval, recurrence_day_of_month,
		       specific_dates, start_date, end_date, created_at, updated_at
		FROM task_recurrences
		WHERE task_id = $1
	`

	row := r.pool.QueryRow(ctx, query, taskID)
	return scanRecurrence(row)
}

func (r *Repository) UpdateRecurrence(ctx context.Context, recurrence *taskdomain.TaskRecurrence) error {
	const query = `
		UPDATE task_recurrences
		SET recurrence_type = $1,
			recurrence_interval = $2,
			recurrence_day_of_month = $3,
			specific_dates = $4,
			start_date = $5,
			end_date = $6,
			updated_at = $7
		WHERE task_id = $8
	`

	_, err := r.pool.Exec(ctx, query,
		recurrence.RecurrenceType,
		recurrence.RecurrenceInterval,
		recurrence.RecurrenceDayOfMonth,
		recurrence.SpecificDates,
		recurrence.StartDate,
		recurrence.EndDate,
		recurrence.UpdatedAt,
		recurrence.TaskID,
	)
	return err
}

func scanTask(scanner interface{ Scan(dest ...any) error }) (*taskdomain.Task, error) {
	var task taskdomain.Task
	var status string

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	task.Status = taskdomain.Status(status)
	return &task, nil
}

func scanTasks(rows pgx.Rows) ([]taskdomain.Task, error) {
	var tasks []taskdomain.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func scanRecurrence(scanner interface{ Scan(dest ...any) error }) (*taskdomain.TaskRecurrence, error) {
	var rec taskdomain.TaskRecurrence
	var recType string
	var dayOfMonth *int
	var endDate *time.Time

	if err := scanner.Scan(
		&rec.TaskID,
		&recType,
		&rec.RecurrenceInterval,
		&dayOfMonth,
		&rec.SpecificDates,
		&rec.StartDate,
		&endDate,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil 
		}
		return nil, err
	}

	rec.RecurrenceType = taskdomain.RecurrenceType(recType)
	rec.RecurrenceDayOfMonth = dayOfMonth
	rec.EndDate = endDate

	return &rec, nil
}
