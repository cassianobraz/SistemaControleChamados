package domain_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestStatus_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status domain.Status
		want   bool
	}{
		{name: "aberto is valid", status: domain.StatusOpen, want: true},
		{name: "em_andamento is valid", status: domain.StatusInProgress, want: true},
		{name: "resolvido is valid", status: domain.StatusResolved, want: true},
		{name: "fechado is valid", status: domain.StatusClosed, want: true},
		{name: "unknown is invalid", status: domain.Status("cancelado"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.status.Valid())
		})
	}
}

func TestStatus_IsOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status domain.Status
		want   bool
	}{
		{name: "aberto counts as open", status: domain.StatusOpen, want: true},
		{name: "em_andamento counts as open", status: domain.StatusInProgress, want: true},
		{name: "resolvido does not count as open", status: domain.StatusResolved, want: false},
		{name: "fechado does not count as open", status: domain.StatusClosed, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.status.IsOpen())
		})
	}
}

func TestStatus_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "aberto", domain.StatusOpen.String())
}
