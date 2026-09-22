package config_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("falls back to defaults when nothing is set", func(t *testing.T) {
		cfg := config.Load()

		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
		assert.Equal(t, "controle_chamados", cfg.MongoDatabase)
	})

	t.Run("reads values from the environment", func(t *testing.T) {
		t.Setenv("PORT", "9090")
		t.Setenv("MONGO_URI", "mongodb://mongo:27017")
		t.Setenv("MONGO_DATABASE", "custom_db")

		cfg := config.Load()

		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, "mongodb://mongo:27017", cfg.MongoURI)
		assert.Equal(t, "custom_db", cfg.MongoDatabase)
	})

	t.Run("ignores an empty override and keeps the default", func(t *testing.T) {
		t.Setenv("PORT", "")

		cfg := config.Load()

		assert.Equal(t, "8080", cfg.Port)
	})
}
