package recurrence

import "time"

type RecurrenceType string

const (
	RecurrenceDaily     RecurrenceType = "daily"
	RecurrenceMonthly   RecurrenceType = "monthly"
	RecurrenceSpecific  RecurrenceType = "specific"
	RecurrenceEvenOdd   RecurrenceType = "even_odd"
)



type RecurrenceSetting struct {
	ID int64 `json:"id"`

	Type           RecurrenceType `json:"type" validate:"required,oneof=daily monthly specific even_odd"`
	IntervalDays   int `json:"interval_days" validate:"omitempty,min=1"`
	IntervalMonths int `json:"interval_months" validate:"omitempty,min=1,max=30"`

	SpecificDay    []time.Time `json:"specific_day,omitempty"`
	EvenOddDays    bool `json:"even_odd_days,omitempty"`

	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecific, RecurrenceEvenOdd:
		return true
	default:
		return false
	}
}