package secrets

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSecretsOpts struct {
	DbUri          string
	DbUsername     string
	DbPassword     string
	DatabaseName   string
	CollectionName string
}

type MongoSecretsManager struct {
	client         *mongo.Client
	collection     *mongo.Collection
	databaseName   string
	collectionName string
}

func mongoConnect(ctx context.Context, client *mongo.Client, databaseName string, collectionName string) (*mongo.Collection, error) {
	err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(databaseName).Collection(collectionName), nil
}

func (s *MongoSecretsManager) PutSecret(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (s *MongoSecretsManager) GetSecret(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

func (s *MongoSecretsManager) Close(ctx context.Context) {
	s.client.Disconnect(ctx)
}

func NewMongoSecretsManager(ctx context.Context, opts MongoSecretsOpts) (SecretsManager, error) {
	var t time.Duration
	client, err := mongo.NewClient(
		mongoOptions.Client().ApplyURI(opts.DbUri),
		&mongoOptions.ClientOptions{
			MaxConnIdleTime: &t,
			Auth: &mongoOptions.Credential{
				Username: opts.DbUsername,
				Password: opts.DbPassword,
			},
		})
	if err != nil {
		return nil, err
	}

	collection, err := mongoConnect(ctx, client, opts.DatabaseName, opts.CollectionName)
	if err != nil {
		return nil, err
	}
	return &MongoSecretsManager{
		client:         client,
		collection:     collection,
		databaseName:   opts.DatabaseName,
		collectionName: opts.CollectionName,
	}, nil
}
