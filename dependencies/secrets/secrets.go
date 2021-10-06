package secrets

import (
	"context"
	"fmt"

	"emarcey/data-vault/common"
)

type SecretsManager interface {
	PutSecret(ctx context.Context, key string, value interface{}) error
	GetSecret(ctx context.Context, key string) (interface{}, error)
	Close(ctx context.Context)
}

type NewSecretsManagerOpts struct {
	managerType string
	mongoOpts   MongoSecretsOpts
}

func NewSecretsManager(ctx context.Context, opts NewSecretsManagerOpts) (SecretsManager, error) {
	switch opts.managerType {
	case "mongodb":
		return NewMongoSecretsManager(ctx, opts.mongoOpts)
	default:
		return nil, common.NewInitializationError("secrets manager", fmt.Sprintf("Unknown secrets manager type %s", opts.managerType))
	}
}
