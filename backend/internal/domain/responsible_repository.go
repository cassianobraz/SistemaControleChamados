package domain

import "context"

type ResponsibleRepository interface {
	Create(ctx context.Context, responsible *Responsible) error
	FindAll(ctx context.Context) ([]*Responsible, error)
	FindByID(ctx context.Context, id string) (*Responsible, error)
	Update(ctx context.Context, responsible *Responsible) error
	Delete(ctx context.Context, id string) error
	CountAll(ctx context.Context) (int64, error)
}
