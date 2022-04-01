package secrets

import (
	"context"

	"emarcey/data-vault/common"
)

type SecretsManager interface {
	CreateSecret(ctx context.Context, secret *common.EncryptedSecret) error
	GetSecret(ctx context.Context, secretId string) (*common.EncryptedSecret, error)
	LogAccess(ctx context.Context, log *common.AccessLog) error
	ListAccessLogs(ctx context.Context, req *common.ListAccessLogsRequest) ([]*common.AccessLog, error)
	Close(ctx context.Context)
}

type SecretsManagerOpts struct {
	ManagerType string           `yaml:"managerType"`
	MongoOpts   MongoSecretsOpts `yaml:"mongoOpts"`
}

func NewSecretsManager(ctx context.Context, opts SecretsManagerOpts) (SecretsManager, error) {
	switch opts.ManagerType {
	case "mongodb":
		return NewMongoSecretsManager(ctx, opts.MongoOpts)
	default:
		return nil, common.NewInitializationError("secrets manager", "Unknown secrets manager type %s", opts.ManagerType)
	}
}
