package mocks

import (
	"context"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ResponsibleRepository struct {
	mock.Mock
}

func (m *ResponsibleRepository) Create(ctx context.Context, responsible *domain.Responsible) error {
	args := m.Called(ctx, responsible)
	return args.Error(0)
}

func (m *ResponsibleRepository) FindAll(ctx context.Context) ([]*domain.Responsible, error) {
	args := m.Called(ctx)

	responsibles, _ := args.Get(0).([]*domain.Responsible)

	return responsibles, args.Error(1)
}

func (m *ResponsibleRepository) FindByID(ctx context.Context, id string) (*domain.Responsible, error) {
	args := m.Called(ctx, id)

	responsible, _ := args.Get(0).(*domain.Responsible)

	return responsible, args.Error(1)
}

func (m *ResponsibleRepository) Update(ctx context.Context, responsible *domain.Responsible) error {
	args := m.Called(ctx, responsible)
	return args.Error(0)
}

func (m *ResponsibleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ResponsibleRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
