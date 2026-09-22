package usecase

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type TicketUseCase interface {
	Create(ctx context.Context, input CreateTicketInput) (*domain.Ticket, error)
	GetByID(ctx context.Context, id string) (*domain.Ticket, error)
	List(ctx context.Context, input ListTicketsInput) (*ListTicketsOutput, error)
	Update(ctx context.Context, id string, input UpdateTicketInput) (*domain.Ticket, error)
	AssignResponsible(ctx context.Context, id string, responsibleID string) (*domain.Ticket, error)
}

type ticketUseCase struct {
	tickets      domain.TicketRepository
	responsibles domain.ResponsibleRepository
}

func NewTicketUseCase(tickets domain.TicketRepository, responsibles domain.ResponsibleRepository) TicketUseCase {
	return &ticketUseCase{tickets: tickets, responsibles: responsibles}
}

func (uc *ticketUseCase) Create(ctx context.Context, input CreateTicketInput) (*domain.Ticket, error) {
	responsible, err := uc.resolveResponsible(ctx, input.ResponsibleID)
	if err != nil {
		return nil, err
	}

	ticket, err := domain.NewTicket(input.Title, input.Description, input.Priority, responsible.ID, responsible.Name)
	if err != nil {
		return nil, err
	}

	if err := uc.tickets.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("criando chamado: %w", err)
	}

	return ticket, nil
}

func (uc *ticketUseCase) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	ticket, err := uc.tickets.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscando chamado %q: %w", id, err)
	}

	return ticket, nil
}

func (uc *ticketUseCase) List(ctx context.Context, input ListTicketsInput) (*ListTicketsOutput, error) {
	filter := domain.TicketFilter{
		Status:        input.Status,
		Priority:      input.Priority,
		ResponsibleID: input.ResponsibleID,
		Search:        input.Search,
		Sort:          input.Sort,
		Page:          normalizePage(input.Page),
		PageSize:      normalizePageSize(input.PageSize),
	}

	if filter.Sort == "" {
		filter.Sort = domain.SortByCreatedAtDesc
	}

	tickets, total, err := uc.tickets.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listando chamados: %w", err)
	}

	return &ListTicketsOutput{
		Tickets:  tickets,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (uc *ticketUseCase) Update(ctx context.Context, id string, input UpdateTicketInput) (*domain.Ticket, error) {
	ticket, err := uc.tickets.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscando chamado %q: %w", id, err)
	}

	if err := ticket.UpdateDetails(input.Title, input.Description, input.Priority, input.Status); err != nil {
		return nil, err
	}

	if err := uc.tickets.Update(ctx, ticket); err != nil {
		return nil, fmt.Errorf("atualizando chamado %q: %w", id, err)
	}

	return ticket, nil
}

func (uc *ticketUseCase) AssignResponsible(ctx context.Context, id string, responsibleID string) (*domain.Ticket, error) {
	ticket, err := uc.tickets.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscando chamado %q: %w", id, err)
	}

	responsible, err := uc.resolveResponsible(ctx, responsibleID)
	if err != nil {
		return nil, err
	}

	ticket.AssignResponsible(responsible.ID, responsible.Name)

	if err := uc.tickets.Update(ctx, ticket); err != nil {
		return nil, fmt.Errorf("atribuindo chamado %q: %w", id, err)
	}

	return ticket, nil
}

func (uc *ticketUseCase) resolveResponsible(ctx context.Context, responsibleID string) (*domain.Responsible, error) {
	if responsibleID != "" {
		responsible, err := uc.responsibles.FindByID(ctx, responsibleID)
		if err != nil {
			return nil, fmt.Errorf("buscando responsável %q: %w", responsibleID, err)
		}

		return responsible, nil
	}

	return uc.pickLeastLoadedResponsible(ctx)
}

func (uc *ticketUseCase) pickLeastLoadedResponsible(ctx context.Context) (*domain.Responsible, error) {
	responsibles, err := uc.responsibles.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listando responsáveis: %w", err)
	}

	if len(responsibles) == 0 {
		return nil, domain.ErrNoResponsibleAvailable
	}

	var (
		candidates []*domain.Responsible
		minLoad    int64 = -1
	)

	for _, responsible := range responsibles {
		if !responsible.Active {
			continue
		}

		load, err := uc.tickets.CountOpenByResponsible(ctx, responsible.ID)
		if err != nil {
			return nil, fmt.Errorf("contando chamados em aberto do responsável %q: %w", responsible.ID, err)
		}

		switch {
		case len(candidates) == 0 || load < minLoad:
			candidates = []*domain.Responsible{responsible}
			minLoad = load
		case load == minLoad:
			candidates = append(candidates, responsible)
		}
	}

	if len(candidates) == 0 {
		return nil, domain.ErrNoResponsibleAvailable
	}

	chosen := candidates[rand.IntN(len(candidates))]

	return chosen, nil
}

func normalizePage(page int) int {
	if page < 1 {
		return defaultPage
	}

	return page
}

func normalizePageSize(pageSize int) int {
	if pageSize < 1 {
		return defaultPageSize
	}

	if pageSize > maxPageSize {
		return maxPageSize
	}

	return pageSize
}
