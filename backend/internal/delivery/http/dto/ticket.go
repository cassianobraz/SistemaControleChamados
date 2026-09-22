package dto

import (
	"time"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

type CreateTicketRequest struct {
	Title         string `json:"title" example:"Impressora do 2º andar não funciona"`
	Description   string `json:"description" example:"A impressora está travando toda vez que alguém tenta imprimir."`
	Priority      string `json:"priority" example:"media"`
	ResponsibleID string `json:"responsible_id,omitempty" example:""`
}

type UpdateTicketRequest struct {
	Title       string `json:"title" example:"Impressora do 2º andar não funciona"`
	Description string `json:"description" example:"A impressora está travando toda vez que alguém tenta imprimir."`
	Priority    string `json:"priority" example:"alta"`
	Status      string `json:"status" example:"em_andamento"`
}

type AssignTicketRequest struct {
	ResponsibleID string `json:"responsible_id,omitempty" example:""`
}

type TicketResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Priority        string    `json:"priority"`
	Status          string    `json:"status"`
	ResponsibleID   string    `json:"responsible_id"`
	ResponsibleName string    `json:"responsible_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewTicketResponse(ticket *domain.Ticket) TicketResponse {
	return TicketResponse{
		ID:              ticket.ID,
		Title:           ticket.Title,
		Description:     ticket.Description,
		Priority:        ticket.Priority.String(),
		Status:          ticket.Status.String(),
		ResponsibleID:   ticket.ResponsibleID,
		ResponsibleName: ticket.ResponsibleName,
		CreatedAt:       ticket.CreatedAt,
		UpdatedAt:       ticket.UpdatedAt,
	}
}

func NewTicketResponseList(tickets []*domain.Ticket) []TicketResponse {
	responses := make([]TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		responses = append(responses, NewTicketResponse(ticket))
	}

	return responses
}

type TicketListResponse struct {
	Data       []TicketResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int64            `json:"total_pages"`
}
