package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/database"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/repository"
)

func setup(t *testing.T) (*repository.TicketRepository, *repository.ResponsibleRepository) {
	t.Helper()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, uri)
	if err != nil {
		t.Skipf("mongodb indisponível em %q, pulando teste de integração: %v", uri, err)
	}

	dbName := "controle_chamados_test_" + bson.NewObjectID().Hex()
	db := client.Database(dbName)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
		_ = client.Disconnect(cleanupCtx)
	})

	return repository.NewTicketRepository(db), repository.NewResponsibleRepository(db)
}

func TestResponsibleRepository_CRUD(t *testing.T) {
	t.Parallel()

	_, responsibleRepo := setup(t)
	ctx := context.Background()

	count, err := responsibleRepo.CountAll(ctx)
	require.NoError(t, err)
	require.Zero(t, count)

	responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
	require.NoError(t, err)
	require.NoError(t, responsibleRepo.Create(ctx, responsible))
	require.NotEmpty(t, responsible.ID)

	found, err := responsibleRepo.FindByID(ctx, responsible.ID)
	require.NoError(t, err)
	require.Equal(t, responsible.Name, found.Name)

	all, err := responsibleRepo.FindAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)

	require.NoError(t, found.UpdateDetails("Ana Paula", "ana.paula@codificar.dev", false))
	require.NoError(t, responsibleRepo.Update(ctx, found))

	updated, err := responsibleRepo.FindByID(ctx, responsible.ID)
	require.NoError(t, err)
	require.Equal(t, "Ana Paula", updated.Name)
	require.False(t, updated.Active)

	require.NoError(t, responsibleRepo.Delete(ctx, responsible.ID))

	_, err = responsibleRepo.FindByID(ctx, responsible.ID)
	require.ErrorIs(t, err, domain.ErrResponsibleNotFound)

	require.ErrorIs(t, responsibleRepo.Delete(ctx, responsible.ID), domain.ErrResponsibleNotFound)

	_, err = responsibleRepo.FindByID(ctx, bson.NewObjectID().Hex())
	require.ErrorIs(t, err, domain.ErrResponsibleNotFound)
}

func TestTicketRepository_CRUD(t *testing.T) {
	t.Parallel()

	ticketRepo, responsibleRepo := setup(t)
	ctx := context.Background()

	responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
	require.NoError(t, err)
	require.NoError(t, responsibleRepo.Create(ctx, responsible))

	ticket, err := domain.NewTicket("Impressora quebrada", "Não liga", domain.PriorityHigh, responsible.ID, responsible.Name)
	require.NoError(t, err)
	require.NoError(t, ticketRepo.Create(ctx, ticket))
	require.NotEmpty(t, ticket.ID)

	found, err := ticketRepo.FindByID(ctx, ticket.ID)
	require.NoError(t, err)
	require.Equal(t, ticket.Title, found.Title)

	openCount, err := ticketRepo.CountOpenByResponsible(ctx, responsible.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), openCount)

	require.NoError(t, found.UpdateDetails(found.Title, found.Description, found.Priority, domain.StatusClosed))
	require.NoError(t, ticketRepo.Update(ctx, found))

	closedCount, err := ticketRepo.CountOpenByResponsible(ctx, responsible.ID)
	require.NoError(t, err)
	require.Zero(t, closedCount)

	tickets, total, err := ticketRepo.FindAll(ctx, domain.TicketFilter{Page: 1, PageSize: 10, Sort: domain.SortByCreatedAtDesc})
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, tickets, 1)

	_, err = ticketRepo.FindByID(ctx, bson.NewObjectID().Hex())
	require.ErrorIs(t, err, domain.ErrTicketNotFound)
}
