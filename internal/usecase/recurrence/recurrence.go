package recurrence

import (
	"context"
	"time"

	taskdomainrec "example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type RepositoryRecurrence interface {
	CreateRecurrence(ctx context.Context, recurrence *taskdomainrec.RecurrenceSetting) (*taskdomainrec.RecurrenceSetting, error)
	GetRecurrenceByID(ctx context.Context, id int64) (*taskdomainrec.RecurrenceSetting, error)
	UpdateRecurrence(ctx context.Context, recurrence *taskdomainrec.RecurrenceSetting) (*taskdomainrec.RecurrenceSetting, error)
	DeleteRecurrence(ctx context.Context, id int64) error
	ListRecurrence(ctx context.Context) ([]taskdomainrec.RecurrenceSetting, error)
	CreateRecurrenceONTask(ctx context.Context, recuurId int64, taskID int64) error
}

type UsecaseRecurrence interface {
	CreateRecurrenceId(ctx context.Context, taskId int64, input CreateInputRecurrence) (*taskdomainrec.RecurrenceSetting, error)
	CreateRecurrence(ctx context.Context, task CreateInputTask, input CreateInputRecurrence) (*taskdomainrec.RecurrenceSetting, error)
	GetByIDRecurrence(ctx context.Context, id int64) (*taskdomainrec.RecurrenceSetting, error)
	UpdateRecurrence(ctx context.Context, id int64, input UpdateInputRecurrence) (*taskdomainrec.RecurrenceSetting, error)
	DeleteRecurrence(ctx context.Context, id int64) error
	ListRecurrence(ctx context.Context) ([]taskdomainrec.RecurrenceSetting, error)
}

type CreateInputRecurrence struct {
	Type           taskdomainrec.RecurrenceType
	IntervalDays   int 
	IntervalMonths int 
	SpecificDay    []time.Time 
	EvenOddDays    bool 

	StartDate time.Time 
	EndDate   time.Time 
}

type UpdateInputRecurrence struct {
	Type           taskdomainrec.RecurrenceType
	IntervalDays   int 
	IntervalMonths int 
	SpecificDay    []time.Time 
	EvenOddDays    bool 

	StartDate time.Time 
	EndDate   time.Time 
}

type CreateInputTask struct {
	Title       string
	Description string
	Status      taskdomain.Status
}

type UpdateInputTask struct {
	Title       string
	Description string
	Status      taskdomain.Status
}
