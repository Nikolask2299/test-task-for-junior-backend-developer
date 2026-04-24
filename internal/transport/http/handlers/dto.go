package handlers

import (
	"time"

	"example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}


type taskReccurrentMutation struct {
	Task taskMutationDTO `json:"task"`
	Recurrence recurrenceMutationDTO `json:"recurrence"`
}

type recurrenceMutationDTO struct {
	Type           recurrence.RecurrenceType `json:"type" validate:"required,oneof=daily monthly specific even_odd"`
	IntervalDays   int `json:"interval_days" validate:"omitempty,min=1"`
	IntervalMonths int `json:"interval_months" validate:"omitempty,min=1,max=30"`
	SpecificDay    []time.Time `json:"specific_day,omitempty" validate:"omitempty,unique"`
	EvenOddDays    bool `json:"even_odd_days,omitempty"`

	StartDate time.Time `json:"start_date,omitempty" validate:"omitempty,datetime"`
	EndDate   time.Time `json:"end_date,omitempty" validate:"omitempty,datetime"`
}

type reccurrentDTO struct {
	ID int64 `json:"id"`

	Type           recurrence.RecurrenceType `json:"type" validate:"required,oneof=daily monthly specific even_odd"`
	IntervalDays   int `json:"interval_days" validate:"omitempty,min=1"`
	IntervalMonths int `json:"interval_months" validate:"omitempty,min=1,max=30"`
	SpecificDay    []time.Time `json:"specific_day,omitempty"`
	EvenOddDays    bool `json:"even_odd_days,omitempty"`

	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

func newReccurrentDTO(task *recurrence.RecurrenceSetting) reccurrentDTO {
	return reccurrentDTO{
		ID: task.ID,
		Type: task.Type,
		IntervalDays: task.IntervalDays,
		IntervalMonths: task.IntervalMonths,
		SpecificDay: task.SpecificDay,
		EvenOddDays: task.EvenOddDays,
		StartDate: task.StartDate,
		EndDate: task.EndDate,
	}
}

	
