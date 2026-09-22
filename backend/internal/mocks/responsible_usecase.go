package mocks

import (
	"context"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// ResponsibleUseCase is a testify mock implementing usecase.ResponsibleUseCase.
type ResponsibleUseCase struct {
	mock.Mock
}

func (m *ResponsibleUseCase) List(ctx context.Context) ([]*usecase.ResponsibleWithWorkload, error) {
	args := m.Called(ctx)

	result, _ := args.Get(0).([]*usecase.ResponsibleWithWorkload)

	return result, args.Error(1)
}

func (m *ResponsibleUseCase) Create(ctx context.Context, input usecase.CreateResponsibleInput) (*usecase.ResponsibleWithWorkload, error) {
	args := m.Called(ctx, input)

	result, _ := args.Get(0).(*usecase.ResponsibleWithWorkload)

	return result, args.Error(1)
}

func (m *ResponsibleUseCase) Update(ctx context.Context, id string, input usecase.UpdateResponsibleInput) (*usecase.ResponsibleWithWorkload, error) {
	args := m.Called(ctx, id, input)

	result, _ := args.Get(0).(*usecase.ResponsibleWithWorkload)

	return result, args.Error(1)
}

func (m *ResponsibleUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
