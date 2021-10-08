package secrets

import (
	"context"
	"fmt"

	"emarcey/data-vault/common"
)

type Secret struct {
	Id         string      `json:"_id" bson:"_id"`
	TableName  string      `json:"table_name" bson:"table_name"`
	ColumnName string      `json:"column_name" bson:"column_name"`
	RowId      string      `json:"row_id" bson:"row_id"`
	IdHash     string      `json:"id_hash" bson:"id_hash"`
	Key        interface{} `json:"key" bson:"key"`
	Iv         interface{} `json:"iv" bson:"iv"`
}

func makeSecretId(tableName string, rowId string, columnName string, idHash string) string {
	return fmt.Sprintf("%s|||%s|||%s|||%s", tableName, rowId, columnName, idHash)
}

func NewSecret(tableName string, rowId string, columnName string, idHash string, key string, iv string) *Secret {
	return &Secret{
		Id:         makeSecretId(tableName, rowId, columnName, idHash),
		TableName:  tableName,
		ColumnName: columnName,
		RowId:      rowId,
		IdHash:     idHash,
		Key:        key,
		Iv:         iv,
	}
}

type SecretsManager interface {
	GetOrPutSecret(ctx context.Context, secret *Secret) (*Secret, error)
	Close(ctx context.Context)
}

type SecretsManagerOpts struct {
	ManagerType string
	MongoOpts   MongoSecretsOpts
}

func NewSecretsManager(ctx context.Context, opts SecretsManagerOpts) (SecretsManager, error) {
	switch opts.ManagerType {
	case "mongodb":
		return NewMongoSecretsManager(ctx, opts.MongoOpts)
	default:
		return nil, common.NewInitializationError("secrets manager", fmt.Sprintf("Unknown secrets manager type %s", opts.ManagerType))
	}
}
