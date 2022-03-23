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
	DbUsername            string `yaml:"dbUsername"`
	DbPassword            string `yaml:"dbPassword"`
	ClusterName           string `yaml:"clusterName"`
	DatabaseName          string `yaml:"databaseName"`
	SecretsCollectionName string `yaml:"secretsCollectionName"`
	LogCollectionName     string `yaml:"logCollectionName"`
}

type MongoSecretsManager struct {
	client                *mongo.Client
	secretsCollection     *mongo.Collection
	logCollection         *mongo.Collection
	databaseName          string
	secretsCollectionName string
	logCollectionName     string
}

func (s *MongoSecretsManager) reconnect(ctx context.Context) error {
	err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	s.secretsCollection = s.client.Database(s.databaseName).Collection(s.secretsCollectionName)
	s.logCollection = s.client.Database(s.databaseName).Collection(s.logCollectionName)
	return nil
}

func (s *MongoSecretsManager) GetSecret(ctx context.Context, secretId string) (*common.EncryptedSecret, error) {
	result := s.secretsCollection.FindOne(ctx, bson.M{"_id": secretId})
	if result == nil {
		return nil, common.NewMongoGetSecretError("FindOne for secret %s returned nil.", secretId)
	}
	err := result.Err()
	if err != nil {
		return nil, err
	}

	var val common.EncryptedSecret
	err = result.Decode(&val)
	if err != nil {
		return nil, common.NewMongoGetSecretError("Decode for secret %s, with raw value %v, returned error: %v.", secretId, result, err)
	}
	return &val, nil
}

func (s *MongoSecretsManager) CreateSecret(ctx context.Context, secret *common.EncryptedSecret) error {
	_, err := s.secretsCollection.InsertOne(ctx, secret)
	if err != nil {
		return common.NewMongoError("CreateSecret", "Error inserting secret, %s, received error, %v", secret.Id, err)
	}
	return nil
}

func (s *MongoSecretsManager) Close(ctx context.Context) {
	s.client.Disconnect(ctx)
}

func (s *MongoSecretsManager) LogAccess(ctx context.Context, log *common.AccessLog) error {
	_, err := s.logCollection.InsertOne(ctx, log)
	if err != nil {
		return common.NewMongoError("LogAccess", "Error inserting log, %+v, received error, %v", log, err)
	}
	return nil
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
		client:                client,
		secretsCollection:     nil,
		logCollection:         nil,
		databaseName:          opts.DatabaseName,
		secretsCollectionName: opts.SecretsCollectionName,
		logCollectionName:     opts.LogCollectionName,
	}
	err = secretsManager.reconnect(ctx)
	if err != nil {
		return nil, err
	}
	return secretsManager, nil
}
