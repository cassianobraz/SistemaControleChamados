package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
)

type TicketHandler struct {
	tickets usecase.TicketUseCase
}

func NewTicketHandler(tickets usecase.TicketUseCase) *TicketHandler {
	return &TicketHandler{tickets: tickets}
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTicketRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ticket, err := h.tickets.Create(r.Context(), usecase.CreateTicketInput{
		Title:         req.Title,
		Description:   req.Description,
		Priority:      domain.Priority(req.Priority),
		ResponsibleID: req.ResponsibleID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ticket, err := h.tickets.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	input := parseListTicketsInput(r)

	output, err := h.tickets.List(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	totalPages := int64(0)
	if output.Total > 0 {
		totalPages = (output.Total + int64(output.PageSize) - 1) / int64(output.PageSize)
	}

	writeJSON(w, http.StatusOK, dto.TicketListResponse{
		Data:       dto.NewTicketResponseList(output.Tickets),
		Total:      output.Total,
		Page:       output.Page,
		PageSize:   output.PageSize,
		TotalPages: totalPages,
	})
}

func parseListTicketsInput(r *http.Request) usecase.ListTicketsInput {
	query := r.URL.Query()

	input := usecase.ListTicketsInput{
		Sort:     domain.SortOrder(query.Get("sort")),
		Page:     atoiOrZero(query.Get("page")),
		PageSize: atoiOrZero(query.Get("page_size")),
	}

	if status := query.Get("status"); status != "" {
		s := domain.Status(status)
		input.Status = &s
	}

	if priority := query.Get("priority"); priority != "" {
		p := domain.Priority(priority)
		input.Priority = &p
	}

	if responsibleID := query.Get("responsible_id"); responsibleID != "" {
		input.ResponsibleID = &responsibleID
	}

	if search := query.Get("search"); search != "" {
		input.Search = &search
	}

	return input
}

func atoiOrZero(value string) int {
	if value == "" {
		return 0
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return n
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateTicketRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ticket, err := h.tickets.Update(r.Context(), id, usecase.UpdateTicketInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    domain.Priority(req.Priority),
		Status:      domain.Status(req.Status),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.AssignTicketRequest
	if r.ContentLength > 0 {
		if err := decodeJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
	}

	ticket, err := h.tickets.AssignResponsible(r.Context(), id, req.ResponsibleID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.NewTicketResponse(ticket))
}
