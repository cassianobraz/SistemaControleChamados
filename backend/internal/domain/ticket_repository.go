package domain

import "context"

type SortOrder string

const (
	SortByCreatedAtDesc SortOrder = "created_at_desc"
	SortByCreatedAtAsc  SortOrder = "created_at_asc"
	SortByPriorityDesc  SortOrder = "priority_desc"
)

type TicketFilter struct {
	Status        *Status
	Priority      *Priority
	ResponsibleID *string
	Search        *string
	Sort          SortOrder
	Page          int
	PageSize      int
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	FindByID(ctx context.Context, id string) (*Ticket, error)
	FindAll(ctx context.Context, filter TicketFilter) ([]*Ticket, int64, error)
	Update(ctx context.Context, ticket *Ticket) error
	CountOpenByResponsible(ctx context.Context, responsibleID string) (int64, error)
}
