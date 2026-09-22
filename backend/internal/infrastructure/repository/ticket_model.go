package repository

import (
	"time"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

type ticketModel struct {
	ID              string    `bson:"_id"`
	Title           string    `bson:"title"`
	Description     string    `bson:"description"`
	Priority        string    `bson:"priority"`
	PriorityRank    int       `bson:"priority_rank"`
	Status          string    `bson:"status"`
	ResponsibleID   string    `bson:"responsible_id"`
	ResponsibleName string    `bson:"responsible_name"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

func priorityRank(priority domain.Priority) int {
	switch priority {
	case domain.PriorityHigh:
		return 3
	case domain.PriorityMedium:
		return 2
	case domain.PriorityLow:
		return 1
	default:
		return 0
	}
}

func newTicketModel(ticket *domain.Ticket) ticketModel {
	return ticketModel{
		ID:              ticket.ID,
		Title:           ticket.Title,
		Description:     ticket.Description,
		Priority:        ticket.Priority.String(),
		PriorityRank:    priorityRank(ticket.Priority),
		Status:          ticket.Status.String(),
		ResponsibleID:   ticket.ResponsibleID,
		ResponsibleName: ticket.ResponsibleName,
		CreatedAt:       ticket.CreatedAt,
		UpdatedAt:       ticket.UpdatedAt,
	}
}

func (m ticketModel) toDomain() *domain.Ticket {
	return &domain.Ticket{
		ID:              m.ID,
		Title:           m.Title,
		Description:     m.Description,
		Priority:        domain.Priority(m.Priority),
		Status:          domain.Status(m.Status),
		ResponsibleID:   m.ResponsibleID,
		ResponsibleName: m.ResponsibleName,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
