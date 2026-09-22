package repository

import (
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

type responsibleModel struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Email  string `bson:"email"`
	Active bool   `bson:"active"`
}

func newResponsibleModel(responsible *domain.Responsible) responsibleModel {
	return responsibleModel{
		ID:     responsible.ID,
		Name:   responsible.Name,
		Email:  responsible.Email,
		Active: responsible.Active,
	}
}

func (m responsibleModel) toDomain() *domain.Responsible {
	return &domain.Responsible{
		ID:     m.ID,
		Name:   m.Name,
		Email:  m.Email,
		Active: m.Active,
	}
}
