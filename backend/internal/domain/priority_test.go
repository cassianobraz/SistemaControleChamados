package domain_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPriority_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		priority domain.Priority
		want     bool
	}{
		{name: "baixa is valid", priority: domain.PriorityLow, want: true},
		{name: "media is valid", priority: domain.PriorityMedium, want: true},
		{name: "alta is valid", priority: domain.PriorityHigh, want: true},
		{name: "empty is invalid", priority: domain.Priority(""), want: false},
		{name: "unknown is invalid", priority: domain.Priority("urgente"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.priority.Valid())
		})
	}
}

func TestPriority_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "alta", domain.PriorityHigh.String())
}
