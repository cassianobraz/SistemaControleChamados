package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode json response", "error", err)
	}
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return errors.New("corpo da requisição inválido")
	}

	return nil
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrTicketNotFound), errors.Is(err, domain.ErrResponsibleNotFound):
		status = http.StatusNotFound
	case isValidationError(err):
		status = http.StatusBadRequest
	}

	if status == http.StatusInternalServerError {
		slog.Error("internal error handling request", "error", err)
		writeJSON(w, status, dto.ErrorResponse{Error: "erro interno do servidor"})

		return
	}

	writeJSON(w, status, dto.ErrorResponse{Error: err.Error()})
}

func isValidationError(err error) bool {
	validationErrors := []error{
		domain.ErrInvalidTitle,
		domain.ErrTitleTooLong,
		domain.ErrInvalidDescription,
		domain.ErrInvalidPriority,
		domain.ErrInvalidStatus,
		domain.ErrInvalidResponsibleName,
		domain.ErrNoResponsibleAvailable,
	}

	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			return true
		}
	}

	return false
}
