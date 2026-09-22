package domain

import "errors"

var (
	ErrInvalidTitle           = errors.New("título é obrigatório")
	ErrTitleTooLong           = errors.New("título deve ter no máximo 200 caracteres")
	ErrInvalidDescription     = errors.New("descrição é obrigatória")
	ErrInvalidPriority        = errors.New("prioridade inválida")
	ErrInvalidStatus          = errors.New("status inválido")
	ErrInvalidResponsibleName = errors.New("nome do responsável é obrigatório")

	ErrTicketNotFound         = errors.New("chamado não encontrado")
	ErrResponsibleNotFound    = errors.New("responsável não encontrado")
	ErrNoResponsibleAvailable = errors.New("nenhum responsável disponível para atribuição automática")
)
