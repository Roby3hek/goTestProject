package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      Status          `json:"status"`
	Recurrence  *TaskRecurrence `json:"recurrence,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type RecurrenceType string

const (
	RecurrenceNone     RecurrenceType = "none"
	RecurrenceDaily    RecurrenceType = "daily"
	RecurrenceMonthly  RecurrenceType = "monthly"
	RecurrenceEven     RecurrenceType = "even"
	RecurrenceOdd      RecurrenceType = "odd"
	RecurrenceSpecific RecurrenceType = "specific"
)

type TaskRecurrence struct {
	TaskID               int64          `json:"task_id"`
	RecurrenceType       RecurrenceType `json:"recurrence_type"`
	RecurrenceInterval   int            `json:"recurrence_interval"`     // для daily
	RecurrenceDayOfMonth *int           `json:"recurrence_day_of_month"` // для monthly (1-31)
	SpecificDates        string         `json:"specific_dates"`          // "2026-04-15,2026-05-20"
	StartDate            time.Time      `json:"start_date"`
	EndDate              *time.Time     `json:"end_date,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceNone, RecurrenceDaily, RecurrenceMonthly, RecurrenceEven, RecurrenceOdd, RecurrenceSpecific:
		return true
	default:
		return false
	}
}
