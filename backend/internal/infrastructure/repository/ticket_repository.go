package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

const ticketsCollectionName = "tickets"

type TicketRepository struct {
	collection *mongo.Collection
}

var _ domain.TicketRepository = (*TicketRepository)(nil)

func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{collection: db.Collection(ticketsCollectionName)}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	model := newTicketModel(ticket)

	if _, err := r.collection.InsertOne(ctx, model); err != nil {
		return fmt.Errorf("inserindo chamado: %w", err)
	}

	return nil
}

func (r *TicketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	var model ticketModel
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrTicketNotFound
		}

		return nil, fmt.Errorf("buscando chamado: %w", err)
	}

	return model.toDomain(), nil
}

func (r *TicketRepository) FindAll(ctx context.Context, filter domain.TicketFilter) ([]*domain.Ticket, int64, error) {
	query := buildTicketQuery(filter)

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("contando chamados: %w", err)
	}

	opts := options.Find().
		SetSkip(int64(filter.Page-1) * int64(filter.PageSize)).
		SetLimit(int64(filter.PageSize)).
		SetSort(sortForTicketFilter(filter.Sort))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("listando chamados: %w", err)
	}
	defer cursor.Close(ctx)

	var models []ticketModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, 0, fmt.Errorf("decodificando chamados: %w", err)
	}

	tickets := make([]*domain.Ticket, 0, len(models))
	for _, model := range models {
		tickets = append(tickets, model.toDomain())
	}

	return tickets, total, nil
}

func (r *TicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	model := newTicketModel(ticket)

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": model.ID}, model)
	if err != nil {
		return fmt.Errorf("atualizando chamado: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrTicketNotFound
	}

	return nil
}

func (r *TicketRepository) CountOpenByResponsible(ctx context.Context, responsibleID string) (int64, error) {
	query := bson.M{
		"responsible_id": responsibleID,
		"status":         bson.M{"$in": []string{domain.StatusOpen.String(), domain.StatusInProgress.String()}},
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("contando chamados em aberto: %w", err)
	}

	return count, nil
}

func buildTicketQuery(filter domain.TicketFilter) bson.M {
	query := bson.M{}

	if filter.Status != nil {
		query["status"] = filter.Status.String()
	}

	if filter.Priority != nil {
		query["priority"] = filter.Priority.String()
	}

	if filter.ResponsibleID != nil {
		query["responsible_id"] = *filter.ResponsibleID
	}

	if filter.Search != nil && *filter.Search != "" {
		regex := bson.M{"$regex": *filter.Search, "$options": "i"}
		query["$or"] = bson.A{
			bson.M{"title": regex},
			bson.M{"description": regex},
		}
	}

	return query
}

func sortForTicketFilter(sort domain.SortOrder) bson.D {
	switch sort {
	case domain.SortByCreatedAtAsc:
		return bson.D{{Key: "created_at", Value: 1}}
	case domain.SortByPriorityDesc:
		return bson.D{{Key: "priority_rank", Value: -1}, {Key: "created_at", Value: -1}}
	case domain.SortByCreatedAtDesc:
		return bson.D{{Key: "created_at", Value: -1}}
	default:
		return bson.D{{Key: "created_at", Value: -1}}
	}
}
