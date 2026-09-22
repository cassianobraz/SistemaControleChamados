package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
)

type ResponsibleHandler struct {
	responsibles usecase.ResponsibleUseCase
}

func NewResponsibleHandler(responsibles usecase.ResponsibleUseCase) *ResponsibleHandler {
	return &ResponsibleHandler{responsibles: responsibles}
}

func (h *ResponsibleHandler) List(w http.ResponseWriter, r *http.Request) {
	responsibles, err := h.responsibles.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.NewResponsibleResponseList(responsibles))
}

func (h *ResponsibleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateResponsibleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	responsible, err := h.responsibles.Create(r.Context(), usecase.CreateResponsibleInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.NewResponsibleResponse(responsible))
}

func (h *ResponsibleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateResponsibleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	responsible, err := h.responsibles.Update(r.Context(), id, usecase.UpdateResponsibleInput{
		Name:   req.Name,
		Email:  req.Email,
		Active: req.Active,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.NewResponsibleResponse(responsible))
}

func (h *ResponsibleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.responsibles.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
