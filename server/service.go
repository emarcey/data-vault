package server

import (
	"context"
	"time"

	"emarcey/data-vault/common"
	"emarcey/data-vault/database"
	"emarcey/data-vault/dependencies"
)

type Service interface {
	Version() string

	// users
	ListUsers(ctx context.Context) ([]*common.User, error)
	GetUser(ctx context.Context, id string) (*common.User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	GetAccessToken(ctx context.Context, user *common.User) (*common.AccessToken, error)

	// keys
	CreateSecret(ctx context.Context, key *CreateSecretRequest) (*common.Secret, error)
	// FetchKey(ctx context.Context, user *common.User, keyName string) (*common.Key, error)
	// UpdateKey(ctx context.Context, user *common.User, key *CreateSecretArgs) (*common.Key, error)
	// DeleteKey(ctx context.Context, user *common.User, keyName string) (*common.Key, error)
}

type service struct {
	version string
	deps    *dependencies.Dependencies
}

func NewService(version string, deps *dependencies.Dependencies) Service {
	return &service{
		version: version,
		deps:    deps,
	}

}

func (s *service) Version() string {
	return s.version
}

func (s *service) ListUsers(ctx context.Context) ([]*common.User, error) {
	return database.ListUsers(ctx, s.deps.Database)
}

func (s *service) GetUser(ctx context.Context, userId string) (*common.User, error) {
	return database.GetUserById(ctx, s.deps.Database, userId)
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	userId := common.GenUuid()
	userSecret := common.GenUuid()
	err := database.CreateUser(ctx, s.deps.Database, userId, req.Name, req.Type, common.HashSha256(userSecret))
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{
		UserId:     userId,
		UserSecret: userSecret,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return err
	}
	err = database.DeprecateLatestAccessToken(ctx, s.deps.Database, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = database.DeleteUser(ctx, s.deps.Database, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *service) GetAccessToken(ctx context.Context, user *common.User) (*common.AccessToken, error) {
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	err = database.DeprecateLatestAccessToken(ctx, tx, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	accessToken := common.GenUuid()
	invalidAt := time.Now().Add(time.Duration(s.deps.ServerConfigs.AccessTokenHours) * time.Hour)
	err = database.CreateAccessToken(ctx, tx, user.Id, common.HashSha256(accessToken), invalidAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &common.AccessToken{
		Id:        accessToken,
		UserId:    user.Id,
		IsLatest:  true,
		InvalidAt: invalidAt,
	}, nil
}

func (s *service) CreateSecret(ctx context.Context, createArgs *CreateSecretRequest) (*common.Secret, error) {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	secretId := common.GenUuid()
	ciphertext, encryptedSecret, err := common.EncryptSecret(secretId, createArgs.Name, common.KEY_SIZE)
	if err != nil {
		return nil, err
	}

	err = s.deps.SecretsManager.CreateSecret(ctx, encryptedSecret)
	if err != nil {
		return nil, err
	}

	secret := &common.Secret{
		Id:          secretId,
		Value:       ciphertext,
		Name:        createArgs.Name,
		Description: createArgs.Description,
		CreatedBy:   user.Id,
		UpdatedBy:   user.Id,
	}
	err = database.CreateSecret(ctx, s.deps.Database, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
