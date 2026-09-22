package seed_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/seed"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResponsibles(t *testing.T) {
	t.Parallel()

	t.Run("creates the default responsibles when the collection is empty", func(t *testing.T) {
		t.Parallel()

		repo := new(mocks.ResponsibleRepository)
		repo.On("CountAll", mock.Anything).Return(int64(0), nil).Once()
		repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Responsible")).Return(nil).Times(3)

		err := seed.Responsibles(context.Background(), repo)

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("does nothing when responsibles already exist", func(t *testing.T) {
		t.Parallel()

		repo := new(mocks.ResponsibleRepository)
		repo.On("CountAll", mock.Anything).Return(int64(3), nil).Once()

		err := seed.Responsibles(context.Background(), repo)

		require.NoError(t, err)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("wraps a failure counting existing responsibles", func(t *testing.T) {
		t.Parallel()

		repo := new(mocks.ResponsibleRepository)
		repo.On("CountAll", mock.Anything).Return(int64(0), errors.New("timeout")).Once()

		err := seed.Responsibles(context.Background(), repo)

		require.Error(t, err)
	})

	t.Run("wraps a failure creating a responsible", func(t *testing.T) {
		t.Parallel()

		repo := new(mocks.ResponsibleRepository)
		repo.On("CountAll", mock.Anything).Return(int64(0), nil).Once()
		repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Responsible")).
			Return(errors.New("timeout")).Once()

		err := seed.Responsibles(context.Background(), repo)

		require.Error(t, err)
	})
}
