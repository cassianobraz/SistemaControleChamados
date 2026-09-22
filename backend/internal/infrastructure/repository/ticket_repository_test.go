package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

func TestBuildTicketQuery(t *testing.T) {
	t.Parallel()

	t.Run("returns an empty filter when nothing is set", func(t *testing.T) {
		t.Parallel()

		query := buildTicketQuery(domain.TicketFilter{})

		assert.Equal(t, bson.M{}, query)
	})

	t.Run("combines every optional filter", func(t *testing.T) {
		t.Parallel()

		status := domain.StatusOpen
		priority := domain.PriorityHigh
		responsibleID := "resp-1"
		search := "impressora"

		query := buildTicketQuery(domain.TicketFilter{
			Status:        &status,
			Priority:      &priority,
			ResponsibleID: &responsibleID,
			Search:        &search,
		})

		assert.Equal(t, "aberto", query["status"])
		assert.Equal(t, "alta", query["priority"])
		assert.Equal(t, "resp-1", query["responsible_id"])
		assert.Contains(t, query, "$or")
	})

	t.Run("ignores an empty search term", func(t *testing.T) {
		t.Parallel()

		empty := ""
		query := buildTicketQuery(domain.TicketFilter{Search: &empty})

		assert.NotContains(t, query, "$or")
	})
}

func TestSortForTicketFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sort domain.SortOrder
		want bson.D
	}{
		{
			name: "created_at_asc",
			sort: domain.SortByCreatedAtAsc,
			want: bson.D{{Key: "created_at", Value: 1}},
		},
		{
			name: "priority_desc",
			sort: domain.SortByPriorityDesc,
			want: bson.D{{Key: "priority_rank", Value: -1}, {Key: "created_at", Value: -1}},
		},
		{
			name: "created_at_desc",
			sort: domain.SortByCreatedAtDesc,
			want: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			name: "unknown falls back to created_at_desc",
			sort: domain.SortOrder("qualquer_coisa"),
			want: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, sortForTicketFilter(tt.sort))
		})
	}
}
