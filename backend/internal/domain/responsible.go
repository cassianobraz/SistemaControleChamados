package domain

import (
	"strings"
	"uuid"
)

type Responsible struct {
	ID     string
	Name   string
	Email  string
	Active bool
}

func NewResponsible(name, email string) (*Responsible, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidResponsibleName
	}

	return &Responsible{
		ID:     uuid.NewV7().String(),
		Name:   name,
		Email:  strings.TrimSpace(email),
		Active: true,
	}, nil
}

func (r *Responsible) UpdateDetails(name, email string, active bool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidResponsibleName
	}

	r.Name = name
	r.Email = strings.TrimSpace(email)
	r.Active = active

	return nil
}
