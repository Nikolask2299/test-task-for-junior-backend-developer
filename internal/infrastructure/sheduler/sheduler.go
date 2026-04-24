package sheduler

import (
	"context"
	"log"
	"time"

	"example.com/taskservice/internal/usecase/recurrence"
	"example.com/taskservice/internal/usecase/task"

	taskdomain "example.com/taskservice/internal/domain/task"
	recudomain "example.com/taskservice/internal/domain/recurrence"
)

type ShedulerRepository interface {
	GetCurrentTaskAndRecurrence(ctx context.Context) ([]*SheduleTicket, error)
}

type SheduleTicket struct {
	Task *taskdomain.Task
	Recurrence *recudomain.RecurrenceSetting
}

type Sheduler struct {
	taskrepo task.Repository
	recurepo recurrence.RepositoryRecurrence
	shedrepo ShedulerRepository
	ctx      context.Context
	timer time.Ticker
}

func NewSheduler(ctx context.Context, taskrepo task.Repository, recurepo recurrence.RepositoryRecurrence, shedrepo ShedulerRepository) *Sheduler {
	return &Sheduler{
		taskrepo: taskrepo,
		recurepo: recurepo,
		shedrepo: shedrepo,
		ctx:      ctx,
		timer:    *time.NewTicker(1 * time.Minute),
	}
}

func (s *Sheduler) Start() error {
	for {
		select {
			case <-s.ctx.Done():
				return nil
			case <-s.timer.C:
				if err := s.ProcessRecurrenceTasks(); err != nil {
					log.Println(err)
					return err
				}
		}
	}
}

func (s *Sheduler)	Stop() {
	s.timer.Stop()
	s.ctx.Done()
}

func (s *Sheduler) ProcessRecurrenceTasks() error {

	tickets, err := s.shedrepo.GetCurrentTaskAndRecurrence(s.ctx)
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		switch ticket.Recurrence.Type {
			case recudomain.RecurrenceDaily:
				if err := s.CreateDailyRecurrenceTask(s.ctx, ticket.Task, ticket.Recurrence); err != nil {
					log.Println(err)
				}
			case recudomain.RecurrenceMonthly:
				if err := s.CreateMonthlyRecurrenceTask(s.ctx, ticket.Task, ticket.Recurrence); err != nil {
					log.Println(err)
				}
			case recudomain.RecurrenceSpecific:
				if err := s.CreateSpecificDayRecurrenceTask(s.ctx, ticket.Task, ticket.Recurrence); err != nil {
					log.Println(err)
				}
			case recudomain.RecurrenceEvenOdd:
				if err := s.CreateEvenOddRecurrenceTask(s.ctx, ticket.Task, ticket.Recurrence); err != nil {
					log.Println(err)
				}
			default:
				log.Println("unknown recurrence type")
		}
	}

	return nil
}

func (s *Sheduler) CreateDailyRecurrenceTask(ctx context.Context, task *taskdomain.Task, recurrence *recudomain.RecurrenceSetting) error {
	now := time.Now().UTC()
	lastTaskCreation := task.CreatedAt

	if now.Sub(lastTaskCreation).Hours()/24 >= float64(recurrence.IntervalDays) {
		_, err := s.taskrepo.Create(ctx, &taskdomain.Task{
			Title:       task.Title,
			Description: task.Description,
			Status:      taskdomain.StatusNew,
		})
		return err
	}
	return nil
}

func (s *Sheduler) CreateMonthlyRecurrenceTask(ctx context.Context, task *taskdomain.Task, recurrence *recudomain.RecurrenceSetting) error {
	now := time.Now().UTC()
	lastTaskCreation := task.CreatedAt

	months := lastTaskCreation.AddDate(0, 0,recurrence.IntervalMonths).UTC()
	if months.After(now) || months.Equal(now) {
		_, err := s.taskrepo.Create(ctx, &taskdomain.Task{
			Title:       task.Title,
			Description: task.Description,
			Status:      taskdomain.StatusNew,
		})
		return err
	}
	return nil
}

func (s *Sheduler) CreateSpecificDayRecurrenceTask(ctx context.Context, task *taskdomain.Task, recurrence *recudomain.RecurrenceSetting) error {
	now := time.Now().UTC()
	today := now.Format("2006-01-02")

	for _, date := range recurrence.SpecificDay {
		if date.Format("2006-01-02") == today {
			_, err := s.taskrepo.Create(ctx, &taskdomain.Task{
				Title:       task.Title,
				Description: task.Description,
				Status:      taskdomain.StatusNew,
			})
			return err
		}
	}
	return nil
}

func (s *Sheduler) CreateEvenOddRecurrenceTask(ctx context.Context, task *taskdomain.Task, recurrence *recudomain.RecurrenceSetting) error {
	now := time.Now().UTC()
	day := now.Day()

	isEven := day%2 == 0
	if (recurrence.EvenOddDays && isEven) || (!recurrence.EvenOddDays && !isEven) {
		_, err := s.taskrepo.Create(ctx, &taskdomain.Task{
			Title:       task.Title,
			Description: task.Description,
			Status:      taskdomain.StatusNew,
		})
		return err
	}
	return nil
}