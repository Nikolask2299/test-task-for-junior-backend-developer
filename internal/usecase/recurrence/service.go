package recurrence

import (
	"context"
	"fmt"
	"time"

	"example.com/taskservice/internal/domain/recurrence"
	taskdomainrec "example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
	"example.com/taskservice/internal/usecase/task"
)

type RecurrenceService struct {
	reporec RepositoryRecurrence
	repotsk task.Repository
	now  func() time.Time
}

func NewRecurrenceService(reporec RepositoryRecurrence, repotsk task.Repository) *RecurrenceService {
	return &RecurrenceService{
		reporec: reporec,
		repotsk: repotsk,
		now:  func() time.Time { return time.Now().UTC()},
	}
}

func (s *RecurrenceService) CreateRecurrenceId(ctx context.Context, taskId int64, input CreateInputRecurrence) (*taskdomainrec.RecurrenceSetting, error) {
	
	if taskId <= 0 {
		return nil, fmt.Errorf("%w: taskId must be positive", ErrInvalidInput)
	}
	
	reccure, err := s.reporec.CreateRecurrence(ctx, &taskdomainrec.RecurrenceSetting{
		Type: input.Type,
		IntervalDays: input.IntervalDays,
		IntervalMonths: input.IntervalMonths,
		SpecificDay: input.SpecificDay,
		EvenOddDays: input.EvenOddDays,

		StartDate: input.StartDate,
		EndDate: input.EndDate,
	})

	if err != nil {
		return nil, err
	}

	err = s.reporec.CreateRecurrenceONTask(ctx, reccure.ID, taskId)
	if err != nil {
		return nil, err
	}

	return reccure, nil
}

func (s *RecurrenceService) CreateRecurrence(ctx context.Context, task CreateInputTask, input CreateInputRecurrence) (*recurrence.RecurrenceSetting, error) {

	taskCreate, err := s.repotsk.Create(ctx, 
		&taskdomain.Task{
			Title: task.Title,
			Description: task.Description,
			Status: task.Status,
		})

	if err != nil {
		return nil, err
	}

	reccureTask, err := s.reporec.CreateRecurrence(ctx,
		&taskdomainrec.RecurrenceSetting{
			Type: input.Type,
			IntervalDays: input.IntervalDays,
			IntervalMonths: input.IntervalMonths,
			SpecificDay: input.SpecificDay,
			EvenOddDays: input.EvenOddDays,

			StartDate: input.StartDate,
			EndDate: input.EndDate,
		})

	if err != nil {
		return nil, err
	}

	err = s.reporec.CreateRecurrenceONTask(ctx, reccureTask.ID, taskCreate.ID)
	if err != nil {
		return nil, err
	}

	return reccureTask, nil
}
	
func (s *RecurrenceService) GetByIDRecurrence(ctx context.Context, id int64) (*taskdomainrec.RecurrenceSetting, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.reporec.GetRecurrenceByID(ctx, id)
}

func (s *RecurrenceService) UpdateRecurrence(ctx context.Context, id int64, input UpdateInputRecurrence) (*taskdomainrec.RecurrenceSetting, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if !input.Type.Valid() {
		return nil, fmt.Errorf("%w: invalid type", ErrInvalidInput)
	}

	model := &taskdomainrec.RecurrenceSetting{
		ID: id,
		Type: input.Type,
		IntervalDays: input.IntervalDays,
		IntervalMonths: input.IntervalMonths,
		SpecificDay: input.SpecificDay,
		EvenOddDays: input.EvenOddDays,

		StartDate: input.StartDate,
		EndDate: input.EndDate,
	}

	if !model.Type.Valid() {
		return nil, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	if !model.StartDate.IsZero() && !model.EndDate.IsZero() && model.StartDate.After(model.EndDate) {
		return nil, fmt.Errorf("%w: start_date cannot be after end_date", ErrInvalidInput)
	}

	if model.Type == recurrence.RecurrenceMonthly && (model.IntervalMonths < 1 || model.IntervalMonths > 30) {
		return nil, fmt.Errorf("%w: interval_months must be between 1 and 30", ErrInvalidInput)
	}

	if model.Type == recurrence.RecurrenceDaily && model.IntervalDays < 1 {
		return nil, fmt.Errorf("%w: interval_days must be at least 1", ErrInvalidInput)
	}

	return s.reporec.UpdateRecurrence(ctx, model)
}
	
func (s *RecurrenceService) DeleteRecurrence(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.reporec.DeleteRecurrence(ctx, id)
}
	
func (s *RecurrenceService) ListRecurrence(ctx context.Context) ([]taskdomainrec.RecurrenceSetting, error) {
	return s.reporec.ListRecurrence(ctx)
}

