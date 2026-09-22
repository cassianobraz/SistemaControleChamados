package usecase

import "github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"

type CreateTicketInput struct {
	Title         string
	Description   string
	Priority      domain.Priority
	ResponsibleID string
}

type UpdateTicketInput struct {
	Title       string
	Description string
	Priority    domain.Priority
	Status      domain.Status
}

type ListTicketsInput struct {
	Status        *domain.Status
	Priority      *domain.Priority
	ResponsibleID *string
	Search        *string
	Sort          domain.SortOrder
	Page          int
	PageSize      int
}

type ListTicketsOutput struct {
	Tickets  []*domain.Ticket
	Total    int64
	Page     int
	PageSize int
}

type ResponsibleWithWorkload struct {
	Responsible *domain.Responsible
	OpenTickets int64
}

type CreateResponsibleInput struct {
	Name  string
	Email string
}

type UpdateResponsibleInput struct {
	Name   string
	Email  string
	Active bool
}
