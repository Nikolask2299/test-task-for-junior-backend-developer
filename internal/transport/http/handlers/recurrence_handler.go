package handlers

import (
	"net/http"

	"example.com/taskservice/internal/usecase/recurrence"
)

type RecurrenceHandler struct {
	shedulercase recurrence.UsecaseRecurrence
}

func NewRecurrenceHandler(shedulercase recurrence.UsecaseRecurrence) *RecurrenceHandler {
	return &RecurrenceHandler{
		shedulercase: shedulercase,
	}
}

func (h *RecurrenceHandler) CreateID(w http.ResponseWriter, r *http.Request) {
	var req recurrenceMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.shedulercase.CreateRecurrenceId(
		r.Context(),
		id,
		recurrence.CreateInputRecurrence{
			Type: req.Type,
			IntervalDays: req.IntervalDays,
			IntervalMonths: req.IntervalMonths,
			SpecificDay: req.SpecificDay,
			EvenOddDays: req.EvenOddDays,
			StartDate: req.StartDate,
			EndDate: req.EndDate,
		})
	
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newReccurrentDTO(created))
}

func (h *RecurrenceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req taskReccurrentMutation;
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.shedulercase.CreateRecurrence(
		r.Context(),
		recurrence.CreateInputTask{
			Title:       req.Task.Title,
			Description: req.Task.Description,
			Status:      req.Task.Status,
		},
		recurrence.CreateInputRecurrence{
			Type: req.Recurrence.Type,
			IntervalDays: req.Recurrence.IntervalDays,
			IntervalMonths: req.Recurrence.IntervalMonths,
			SpecificDay: req.Recurrence.SpecificDay,
			EvenOddDays: req.Recurrence.EvenOddDays,

			StartDate: req.Recurrence.StartDate,
			EndDate: req.Recurrence.EndDate,
		})

	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	
	writeJSON(w, http.StatusCreated, newReccurrentDTO(created))
}

func (h *RecurrenceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	recurrence, err := h.shedulercase.GetByIDRecurrence(r.Context(), id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newReccurrentDTO(recurrence))
}

func (h *RecurrenceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req recurrenceMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := h.shedulercase.UpdateRecurrence(
		r.Context(),
		id,
		recurrence.UpdateInputRecurrence{
			Type: req.Type,
			IntervalDays: req.IntervalDays,
			IntervalMonths: req.IntervalMonths,
			SpecificDay: req.SpecificDay,
			EvenOddDays: req.EvenOddDays,

			StartDate: req.StartDate,
			EndDate: req.EndDate,
		})

	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newReccurrentDTO(updated))
}

func (h *RecurrenceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.shedulercase.DeleteRecurrence(r.Context(), id); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecurrenceHandler) List(w http.ResponseWriter, r *http.Request) {
	recurrences, err := h.shedulercase.ListRecurrence(r.Context())
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	response := make([]reccurrentDTO, 0, len(recurrences))
	for i := range recurrences {
		response = append(response, newReccurrentDTO(&recurrences[i]))
	}

	writeJSON(w, http.StatusOK, response)
}


