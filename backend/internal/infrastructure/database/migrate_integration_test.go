//go:build integration

// Run with:
//
//	MONGO_URI=mongodb://localhost:27017 go test -tags=integration ./internal/infrastructure/database/...
package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/database"
)

func TestMigrateLegacyObjectIDs(t *testing.T) {
	t.Parallel()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, uri)
	if err != nil {
		t.Skipf("mongodb indisponível em %q, pulando teste de integração: %v", uri, err)
	}

	db := client.Database("controle_chamados_test_" + bson.NewObjectID().Hex())
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
		_ = client.Disconnect(cleanupCtx)
	})

	collection := db.Collection("tickets")
	legacyID := bson.NewObjectID()

	_, err = collection.InsertOne(ctx, bson.M{"_id": legacyID, "title": "chamado antigo"})
	require.NoError(t, err)

	require.NoError(t, database.MigrateLegacyObjectIDs(ctx, db, "tickets", "responsibles"))

	var migrated bson.M
	require.NoError(t, collection.FindOne(ctx, bson.M{"_id": legacyID.Hex()}).Decode(&migrated))
	require.Equal(t, "chamado antigo", migrated["title"])

	count, err := collection.CountDocuments(ctx, bson.M{"_id": legacyID})
	require.NoError(t, err)
	require.Zero(t, count, "the original ObjectID-keyed document should have been removed")

	// Running it again must be a no-op: nothing left to migrate.
	require.NoError(t, database.MigrateLegacyObjectIDs(ctx, db, "tickets", "responsibles"))

	total, err := collection.CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
}
