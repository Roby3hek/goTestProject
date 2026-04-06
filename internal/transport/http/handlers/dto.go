package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	// Периодичность
	RecurrenceType       taskdomain.RecurrenceType `json:"recurrence_type,omitempty"`
	RecurrenceInterval   int                       `json:"recurrence_interval,omitempty"`
	RecurrenceDayOfMonth *int                      `json:"recurrence_day_of_month,omitempty"`
	SpecificDates        string                    `json:"specific_dates,omitempty"`
	StartDate            *time.Time                `json:"start_date,omitempty"`
	EndDate              *time.Time                `json:"end_date,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	// Периодичность (опционально)
	Recurrence *taskRecurrenceDTO `json:"recurrence,omitempty"`
}

type taskRecurrenceDTO struct {
	RecurrenceType       taskdomain.RecurrenceType `json:"recurrence_type"`
	RecurrenceInterval   int                       `json:"recurrence_interval"`
	RecurrenceDayOfMonth *int                      `json:"recurrence_day_of_month,omitempty"`
	SpecificDates        string                    `json:"specific_dates,omitempty"`
	StartDate            time.Time                 `json:"start_date"`
	EndDate              *time.Time                `json:"end_date,omitempty"`
}

func newTaskDTO(task *taskdomain.Task, recurrence *taskdomain.TaskRecurrence) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if recurrence != nil && recurrence.RecurrenceType != taskdomain.RecurrenceNone {
		dto.Recurrence = &taskRecurrenceDTO{
			RecurrenceType:       recurrence.RecurrenceType,
			RecurrenceInterval:   recurrence.RecurrenceInterval,
			RecurrenceDayOfMonth: recurrence.RecurrenceDayOfMonth,
			SpecificDates:        recurrence.SpecificDates,
			StartDate:            recurrence.StartDate,
			EndDate:              recurrence.EndDate,
		}
	}

	return dto
}
