package domain

import (
	"strings"
	"time"
	"uuid"
)

const maxTitleLength = 200

type Ticket struct {
	ID              string
	Title           string
	Description     string
	Priority        Priority
	Status          Status
	ResponsibleID   string
	ResponsibleName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewTicket(title, description string, priority Priority, responsibleID, responsibleName string) (*Ticket, error) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if err := validateTitle(title); err != nil {
		return nil, err
	}

	if description == "" {
		return nil, ErrInvalidDescription
	}

	if !priority.Valid() {
		return nil, ErrInvalidPriority
	}

	now := time.Now().UTC()

	return &Ticket{
		ID:              uuid.NewV7().String(),
		Title:           title,
		Description:     description,
		Priority:        priority,
		Status:          StatusOpen,
		ResponsibleID:   responsibleID,
		ResponsibleName: responsibleName,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func validateTitle(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}

	if len(title) > maxTitleLength {
		return ErrTitleTooLong
	}

	return nil
}

func (t *Ticket) UpdateDetails(title, description string, priority Priority, status Status) error {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if err := validateTitle(title); err != nil {
		return err
	}

	if description == "" {
		return ErrInvalidDescription
	}

	if !priority.Valid() {
		return ErrInvalidPriority
	}

	if !status.Valid() {
		return ErrInvalidStatus
	}

	t.Title = title
	t.Description = description
	t.Priority = priority
	t.Status = status
	t.UpdatedAt = time.Now().UTC()

	return nil
}

func (t *Ticket) AssignResponsible(responsibleID, responsibleName string) {
	t.ResponsibleID = responsibleID
	t.ResponsibleName = responsibleName
	t.UpdatedAt = time.Now().UTC()
}

func (t *Ticket) IsOpen() bool {
	return t.Status.IsOpen()
}
