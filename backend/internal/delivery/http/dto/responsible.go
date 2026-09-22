package dto

import "github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"

type CreateResponsibleRequest struct {
	Name  string `json:"name" example:"Ana Souza"`
	Email string `json:"email" example:"ana@codificar.dev"`
}

type UpdateResponsibleRequest struct {
	Name   string `json:"name" example:"Ana Souza"`
	Email  string `json:"email" example:"ana@codificar.dev"`
	Active bool   `json:"active" example:"true"`
}

type ResponsibleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Active      bool   `json:"active"`
	OpenTickets int64  `json:"open_tickets"`
}

func NewResponsibleResponse(item *usecase.ResponsibleWithWorkload) ResponsibleResponse {
	return ResponsibleResponse{
		ID:          item.Responsible.ID,
		Name:        item.Responsible.Name,
		Email:       item.Responsible.Email,
		Active:      item.Responsible.Active,
		OpenTickets: item.OpenTickets,
	}
}

func NewResponsibleResponseList(items []*usecase.ResponsibleWithWorkload) []ResponsibleResponse {
	responses := make([]ResponsibleResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, NewResponsibleResponse(item))
	}

	return responses
}
