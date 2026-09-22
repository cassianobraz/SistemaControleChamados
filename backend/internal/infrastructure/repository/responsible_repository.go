package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

const responsiblesCollectionName = "responsibles"

type ResponsibleRepository struct {
	collection *mongo.Collection
}

var _ domain.ResponsibleRepository = (*ResponsibleRepository)(nil)

func NewResponsibleRepository(db *mongo.Database) *ResponsibleRepository {
	return &ResponsibleRepository{collection: db.Collection(responsiblesCollectionName)}
}

func (r *ResponsibleRepository) Create(ctx context.Context, responsible *domain.Responsible) error {
	model := newResponsibleModel(responsible)

	if _, err := r.collection.InsertOne(ctx, model); err != nil {
		return fmt.Errorf("inserindo responsável: %w", err)
	}

	return nil
}

func (r *ResponsibleRepository) FindAll(ctx context.Context) ([]*domain.Responsible, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("listando responsáveis: %w", err)
	}
	defer cursor.Close(ctx)

	var models []responsibleModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("decodificando responsáveis: %w", err)
	}

	responsibles := make([]*domain.Responsible, 0, len(models))
	for _, model := range models {
		responsibles = append(responsibles, model.toDomain())
	}

	return responsibles, nil
}

func (r *ResponsibleRepository) FindByID(ctx context.Context, id string) (*domain.Responsible, error) {
	var model responsibleModel
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrResponsibleNotFound
		}

		return nil, fmt.Errorf("buscando responsável: %w", err)
	}

	return model.toDomain(), nil
}

func (r *ResponsibleRepository) Update(ctx context.Context, responsible *domain.Responsible) error {
	model := newResponsibleModel(responsible)

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": model.ID}, model)
	if err != nil {
		return fmt.Errorf("atualizando responsável: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrResponsibleNotFound
	}

	return nil
}

func (r *ResponsibleRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("removendo responsável: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrResponsibleNotFound
	}

	return nil
}

func (r *ResponsibleRepository) CountAll(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("contando responsáveis: %w", err)
	}

	return count, nil
}
