package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()

	task := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	if normalized.RecurrenceType != taskdomain.RecurrenceNone {
		recurrence := &taskdomain.TaskRecurrence{
			TaskID:               createdTask.ID,
			RecurrenceType:       normalized.RecurrenceType,
			RecurrenceInterval:   normalized.RecurrenceInterval,
			RecurrenceDayOfMonth: normalized.RecurrenceDayOfMonth,
			SpecificDates:        normalized.SpecificDates,
			StartDate:            normalized.StartDate,
			EndDate:              normalized.EndDate,
			CreatedAt:            now,
			UpdatedAt:            now,
		}

		if err := s.repo.CreateRecurrence(ctx, recurrence); err != nil {
			return nil, fmt.Errorf("failed to create recurrence: %w", err)
		}
	}

	return createdTask, nil
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()

	task := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   now,
	}

	updatedTask, err := s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	if normalized.HasRecurrence {
		recurrence := &taskdomain.TaskRecurrence{
			TaskID:               id,
			RecurrenceType:       normalized.RecurrenceType,
			RecurrenceInterval:   normalized.RecurrenceInterval,
			RecurrenceDayOfMonth: normalized.RecurrenceDayOfMonth,
			SpecificDates:        normalized.SpecificDates,
			StartDate:            normalized.StartDate,
			EndDate:              normalized.EndDate,
			UpdatedAt:            now,
		}

		existing, _ := s.repo.GetRecurrenceByTaskID(ctx, id)
		if existing != nil {
			if err := s.repo.UpdateRecurrence(ctx, recurrence); err != nil {
				return nil, fmt.Errorf("failed to update recurrence: %w", err)
			}
		} else {
			recurrence.CreatedAt = now
			if err := s.repo.CreateRecurrence(ctx, recurrence); err != nil {
				return nil, fmt.Errorf("failed to create recurrence: %w", err)
			}
		}
	}

	return updatedTask, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}
	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validateRecurrence(input.RecurrenceType, input.RecurrenceInterval,
		input.RecurrenceDayOfMonth, input.SpecificDates, input.StartDate, input.EndDate); err != nil {
		return CreateInput{}, err
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.HasRecurrence {
		if err := validateRecurrence(input.RecurrenceType, input.RecurrenceInterval,
			input.RecurrenceDayOfMonth, input.SpecificDates, input.StartDate, input.EndDate); err != nil {
			return UpdateInput{}, err
		}
	}

	return input, nil
}

func validateRecurrence(recType taskdomain.RecurrenceType, interval int, dayOfMonth *int,
	specificDates string, startDate time.Time, endDate *time.Time) error {

	if recType == taskdomain.RecurrenceNone {
		return nil
	}

	if !recType.Valid() {
		return fmt.Errorf("%w: invalid recurrence_type", ErrInvalidInput)
	}

	if startDate.IsZero() {
		return fmt.Errorf("%w: start_date is required for recurring tasks", ErrInvalidInput)
	}

	switch recType {
	case taskdomain.RecurrenceDaily:
		if interval < 1 {
			return fmt.Errorf("%w: recurrence_interval must be >= 1 for daily", ErrInvalidInput)
		}
	case taskdomain.RecurrenceMonthly:
		if dayOfMonth == nil || *dayOfMonth < 1 || *dayOfMonth > 31 {
			return fmt.Errorf("%w: recurrence_day_of_month must be between 1 and 31", ErrInvalidInput)
		}
	case taskdomain.RecurrenceSpecific:
		if strings.TrimSpace(specificDates) == "" {
			return fmt.Errorf("%w: specific_dates is required", ErrInvalidInput)
		}
	}

	if endDate != nil && endDate.Before(startDate) {
		return fmt.Errorf("%w: end_date cannot be before start_date", ErrInvalidInput)
	}

	return nil
}
