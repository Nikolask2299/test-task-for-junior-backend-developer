package postgres

import (
	"context"

	taskdomain "example.com/taskservice/internal/domain/task"
	recurdomain "example.com/taskservice/internal/domain/recurrence"
	"example.com/taskservice/internal/infrastructure/sheduler"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShedulerRepository struct {
	db *pgxpool.Pool
}

func NewShedulerRepository(db *pgxpool.Pool) *ShedulerRepository {
	return &ShedulerRepository{
		db: db,
	}
}


func (r *ShedulerRepository) GetCurrentTaskAndRecurrence(ctx context.Context) ([]*sheduler.SheduleTicket, error) {
	const query = `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at, r.id, r.type, r.interval_days, r.interval_months, r.specific_day, r.even_odd_days, r.start_date, r.end_date
		FROM tasks t
		JOIN (SELECT recurrence_id, MAX(task_id) as task_id FROM recurrence_on_task 
		GROUP BY recurrence_id) rt ON t.id = rt.task_id
		JOIN recurrences r ON r.id = rt.recurrence_id
		WHERE r.start_date <= NOW() AND (r.end_date IS NULL OR r.end_date >= NOW())
		`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tickets := make([]*sheduler.SheduleTicket, 0)
	for rows.Next() {
		ticket, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

type ticketScanner interface {
	Scan(dest ...any) error
}

func scanTicket(scanner taskScanner) (*sheduler.SheduleTicket, error) {
	var (
		ticket   sheduler.SheduleTicket
		status string
		typeString string
	)

	if err := scanner.Scan(
		&ticket.Task.ID,
		&ticket.Task.Title,
		&ticket.Task.Description,
		status,
		&ticket.Task.CreatedAt,
		&ticket.Task.UpdatedAt,
		&ticket.Recurrence.ID,
		typeString,
		&ticket.Recurrence.IntervalDays,
		&ticket.Recurrence.IntervalMonths,
		&ticket.Recurrence.SpecificDay,
		&ticket.Recurrence.EvenOddDays,
		&ticket.Recurrence.StartDate,
		&ticket.Recurrence.EndDate,
	); err != nil {
		return nil, err
	}

	ticket.Task.Status = taskdomain.Status(status)
	ticket.Recurrence.Type = recurdomain.RecurrenceType(typeString)

	return &ticket, nil
}