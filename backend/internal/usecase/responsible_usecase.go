package usecase

import (
	"context"
	"fmt"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

type ResponsibleUseCase interface {
	List(ctx context.Context) ([]*ResponsibleWithWorkload, error)
	Create(ctx context.Context, input CreateResponsibleInput) (*ResponsibleWithWorkload, error)
	Update(ctx context.Context, id string, input UpdateResponsibleInput) (*ResponsibleWithWorkload, error)
	Delete(ctx context.Context, id string) error
}

type responsibleUseCase struct {
	responsibles domain.ResponsibleRepository
	tickets      domain.TicketRepository
}

func NewResponsibleUseCase(responsibles domain.ResponsibleRepository, tickets domain.TicketRepository) ResponsibleUseCase {
	return &responsibleUseCase{responsibles: responsibles, tickets: tickets}
}

func (uc *responsibleUseCase) List(ctx context.Context) ([]*ResponsibleWithWorkload, error) {
	responsibles, err := uc.responsibles.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listando responsáveis: %w", err)
	}

	result := make([]*ResponsibleWithWorkload, 0, len(responsibles))

	for _, responsible := range responsibles {
		openTickets, err := uc.tickets.CountOpenByResponsible(ctx, responsible.ID)
		if err != nil {
			return nil, fmt.Errorf("contando chamados em aberto do responsável %q: %w", responsible.ID, err)
		}

		result = append(result, &ResponsibleWithWorkload{
			Responsible: responsible,
			OpenTickets: openTickets,
		})
	}

	return result, nil
}

func (uc *responsibleUseCase) Create(ctx context.Context, input CreateResponsibleInput) (*ResponsibleWithWorkload, error) {
	responsible, err := domain.NewResponsible(input.Name, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.responsibles.Create(ctx, responsible); err != nil {
		return nil, fmt.Errorf("criando responsável: %w", err)
	}

	return &ResponsibleWithWorkload{Responsible: responsible, OpenTickets: 0}, nil
}

func (uc *responsibleUseCase) Update(ctx context.Context, id string, input UpdateResponsibleInput) (*ResponsibleWithWorkload, error) {
	responsible, err := uc.responsibles.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscando responsável %q: %w", id, err)
	}

	if err := responsible.UpdateDetails(input.Name, input.Email, input.Active); err != nil {
		return nil, err
	}

	if err := uc.responsibles.Update(ctx, responsible); err != nil {
		return nil, fmt.Errorf("atualizando responsável %q: %w", id, err)
	}

	openTickets, err := uc.tickets.CountOpenByResponsible(ctx, responsible.ID)
	if err != nil {
		return nil, fmt.Errorf("contando chamados em aberto do responsável %q: %w", responsible.ID, err)
	}

	return &ResponsibleWithWorkload{Responsible: responsible, OpenTickets: openTickets}, nil
}

func (uc *responsibleUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.responsibles.Delete(ctx, id); err != nil {
		return fmt.Errorf("removendo responsável %q: %w", id, err)
	}

	return nil
}
