package postgres

import (
	"context"
	"database/sql"
	"errors"

	"example.com/taskservice/internal/domain/recurrence"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecurrenceRepository struct {
	pool *pgxpool.Pool
}

func NewRecurrenceRepository(pool *pgxpool.Pool) *RecurrenceRepository {
	return &RecurrenceRepository{
		pool: pool,
	}
}

func (r *RecurrenceRepository) CreateRecurrenceONTask(ctx context.Context, recuurId int64, taskID int64) error {
	const query = `
		INSERT INTO recurrence_on_task (task_id, recurrence_id)
		VALUES ($1, $2)
		`

	err := r.pool.QueryRow(ctx, query, taskID, recuurId).Scan()
	if err != nil {
		return err
	}

	return nil
}

func (r *RecurrenceRepository) CreateRecurrence(ctx context.Context, recurrenceSetting *recurrence.RecurrenceSetting) (*recurrence.RecurrenceSetting, error) {
	const query = `
		INSERT INTO recurrences (type, interval_days, interval_months, specific_day, even_odd_days, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, type, interval_days, interval_months, specific_day, even_odd_days, start_date, end_date`
	
	row := r.pool.QueryRow(ctx, query, 
		recurrenceSetting.Type,
		recurrenceSetting.IntervalDays,
		recurrenceSetting.IntervalMonths,
		recurrenceSetting.SpecificDay,
		recurrenceSetting.EvenOddDays,
		recurrenceSetting.StartDate,
		recurrenceSetting.EndDate)
	
	created, err := scanRecurrence(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *RecurrenceRepository) GetRecurrenceByID(ctx context.Context, id int64) (*recurrence.RecurrenceSetting, error) {
	const query = `
		SELECT id, type, interval_days, interval_months, specific_day, even_odd_days, start_date, end_date
		FROM recurrences
		WHERE id = $1`
	
	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanRecurrence(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, recurrence.ErrNotFound
		}
		return nil, err
	}

	return found, nil
}

func (r *RecurrenceRepository) UpdateRecurrence(ctx context.Context, recurrenceSetting *recurrence.RecurrenceSetting)(*recurrence.RecurrenceSetting, error) {
	const query = `
		UPDATE recurrences
		SET type = $1, interval_days = $2, interval_months = $3, specific_day = $4, even_odd_days = $5, start_date = $6, end_date = $7
		WHERE id = $8
		RETURNING id, type, interval_days, interval_months, specific_day, even_odd_days, start_date, end_date`

	row := r.pool.QueryRow(ctx, query,
		recurrenceSetting.Type,
		recurrenceSetting.IntervalDays,
		recurrenceSetting.IntervalMonths,
		recurrenceSetting.SpecificDay,
		recurrenceSetting.EvenOddDays,
		recurrenceSetting.StartDate,
		recurrenceSetting.EndDate,
		recurrenceSetting.ID)
	
	update, err := scanRecurrence(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, recurrence.ErrNotFound
		}

		return nil, err
	}

	return update, nil
}

func (r *RecurrenceRepository) DeleteRecurrence(ctx context.Context, id int64) error {
	const query = `DELETE FROM recurrences WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return recurrence.ErrNotFound
	}

	return nil
}

func (r *RecurrenceRepository) ListRecurrence(ctx context.Context) ([]recurrence.RecurrenceSetting, error) {
	const query = `
		SELECT id, type, interval_days, interval_months, specific_day, even_odd_days, start_date, end_date
		FROM recurrences
		ORDER BY id DESC	
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recurrence := make([]recurrence.RecurrenceSetting, 0)
	for rows.Next() {
		rec, err := scanRecurrence(rows)
		if err != nil {
			return nil, err
		}

		recurrence = append(recurrence, *rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recurrence, nil
}

type recurrenceScanner interface {
	Scan(dest ...any) error
}

func scanRecurrence(scanner recurrenceScanner) (*recurrence.RecurrenceSetting, error) {
	var (
		recurrenceSetting recurrence.RecurrenceSetting
		typeString               string
	)

	if err := scanner.Scan(
		&recurrenceSetting.ID,
		&typeString,
		&recurrenceSetting.IntervalDays,
		&recurrenceSetting.IntervalMonths,
		&recurrenceSetting.SpecificDay,
		&recurrenceSetting.EvenOddDays,
		&recurrenceSetting.StartDate,
		&recurrenceSetting.EndDate,
	); err != nil {
		return nil, err
	}

	recurrenceSetting.Type = recurrence.RecurrenceType(typeString)

	return &recurrenceSetting, nil
}

