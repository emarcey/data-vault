package secrets

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	mongoWriteConcern "go.mongodb.org/mongo-driver/mongo/writeconcern"

	"emarcey/data-vault/common"
)

type MongoSecretsOpts struct {
	DbUsername     string `yaml:"dbUsername"`
	DbPassword     string `yaml:"dbPassword"`
	ClusterName    string `yaml:"clusterName"`
	DatabaseName   string `yaml:"databaseName"`
	CollectionName string `yaml:"collectionName"`
}

type MongoSecretsManager struct {
	client         *mongo.Client
	collection     *mongo.Collection
	databaseName   string
	collectionName string
}

func (s *MongoSecretsManager) reconnect(ctx context.Context) error {
	err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	s.collection = s.client.Database(s.databaseName).Collection(s.collectionName)
	return nil
}

func (s *MongoSecretsManager) GetOrPutSecret(ctx context.Context, secret *Secret) (*Secret, error) {
	result := s.collection.FindOne(ctx, bson.M{"_id": secret.Id})
	if result == nil {
		return nil, common.NewMongoGetOrPutSecretError("FindOne for secret %s returned nil.", secret.Id)
	}
	err := result.Err()
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			return nil, common.NewMongoGetOrPutSecretError("FindOne for secret %s returned error: %v.", secret.Id, err)
		}
		result, err := s.collection.InsertOne(ctx, secret)
		if err != nil {
			return nil, common.NewMongoGetOrPutSecretError("Error inserting secret, %s, received error, %v", secret.Id, err)
		}
		if result == nil {
			return nil, common.NewMongoGetOrPutSecretError("Nil result inserting secret, %s", secret.Id)
		}
		return secret, nil
	}

	var val Secret
	err = result.Decode(&val)
	if err != nil {
		return nil, common.NewMongoGetOrPutSecretError("Decode for secret %s, with raw value %v, returned error: %v.", secret.Id, result, err)
	}
	return &val, nil
}

func (s *MongoSecretsManager) Close(ctx context.Context) {
	s.client.Disconnect(ctx)
}

func NewMongoSecretsManager(ctx context.Context, opts MongoSecretsOpts) (SecretsManager, error) {
	var t time.Duration
	retryWrites := true
	client, err := mongo.NewClient(
		mongoOptions.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", opts.DbUsername, opts.DbPassword, opts.ClusterName, opts.DatabaseName)),
		&mongoOptions.ClientOptions{
			MaxConnIdleTime: &t,
			RetryWrites:     &retryWrites,
			WriteConcern:    mongoWriteConcern.New(mongoWriteConcern.WMajority()),
		})
	if err != nil {
		return nil, err
	}

	secretsManager := &MongoSecretsManager{
		client:         client,
		collection:     nil,
		databaseName:   opts.DatabaseName,
		collectionName: opts.CollectionName,
	}
	err = secretsManager.reconnect(ctx)
	if err != nil {
		return nil, err
	}
	return secretsManager, nil
}
