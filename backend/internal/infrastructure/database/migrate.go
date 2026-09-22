package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MigrateLegacyObjectIDs(ctx context.Context, db *mongo.Database, collections ...string) error {
	for _, name := range collections {
		if err := migrateCollectionIDs(ctx, db.Collection(name)); err != nil {
			return fmt.Errorf("migrando ids legados de %q: %w", name, err)
		}
	}

	return nil
}

func migrateCollectionIDs(ctx context.Context, collection *mongo.Collection) error {
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$type": "objectId"}})
	if err != nil {
		return fmt.Errorf("buscando documentos legados: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := bson.Unmarshal(cursor.Current, &doc); err != nil {
			return fmt.Errorf("decodificando documento legado: %w", err)
		}

		oid, ok := doc["_id"].(bson.ObjectID)
		if !ok {
			continue
		}

		doc["_id"] = oid.Hex()

		if _, err := collection.InsertOne(ctx, doc); err != nil {
			return fmt.Errorf("inserindo documento migrado %q: %w", oid.Hex(), err)
		}

		if _, err := collection.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
			return fmt.Errorf("removendo documento legado %q: %w", oid.Hex(), err)
		}
	}

	return cursor.Err()
}
