package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetBSONOptions(&options.BSONOptions{ObjectIDAsHexString: true})

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("conectando ao mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("verificando conexão com o mongodb: %w", err)
	}

	return client, nil
}
