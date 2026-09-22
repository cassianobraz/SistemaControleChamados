package seed

import (
	"context"
	"fmt"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

var defaultResponsibles = []struct {
	Name  string
	Email string
}{
	{Name: "Ana Souza", Email: "ana.souza@codificar.dev"},
	{Name: "Bruno Lima", Email: "bruno.lima@codificar.dev"},
	{Name: "Carla Mendes", Email: "carla.mendes@codificar.dev"},
}

func Responsibles(ctx context.Context, repo domain.ResponsibleRepository) error {
	count, err := repo.CountAll(ctx)
	if err != nil {
		return fmt.Errorf("verificando responsáveis existentes: %w", err)
	}

	if count > 0 {
		return nil
	}

	for _, seedResponsible := range defaultResponsibles {
		responsible, err := domain.NewResponsible(seedResponsible.Name, seedResponsible.Email)
		if err != nil {
			return fmt.Errorf("construindo responsável %q: %w", seedResponsible.Name, err)
		}

		if err := repo.Create(ctx, responsible); err != nil {
			return fmt.Errorf("criando responsável %q: %w", seedResponsible.Name, err)
		}
	}

	return nil
}
