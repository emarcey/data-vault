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
	GetUser(ctx context.Context, userId string) (*common.User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	RotateUserSecret(ctx context.Context) (*CreateUserResponse, error)
	DeleteUser(ctx context.Context, userId string) error
	GetAccessToken(ctx context.Context) (*common.AccessToken, error)

	// secrets
	CreateSecret(ctx context.Context, key *CreateSecretRequest) (*common.Secret, error)
	GetSecret(ctx context.Context, secretName string) (*common.Secret, error)
	DeleteSecret(ctx context.Context, secretName string) error
	GrantPermission(ctx context.Context, req *SecretPermissionRequest) error
	RevokePermission(ctx context.Context, req *SecretPermissionRequest) error
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
	user, err := database.CreateUser(ctx, s.deps.Database, userId, req.Name, req.Type, common.HashSha256(userSecret))
	if err != nil {
		return nil, err
	}
	s.deps.AuthUsers.Add(userId, user)

	return &CreateUserResponse{
		UserId:     userId,
		UserSecret: userSecret,
		StatusCode: 201,
	}, nil
}

func (s *service) RotateUserSecret(ctx context.Context) (*CreateUserResponse, error) {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userSecret := common.GenUuid()
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	tokenId, err := database.DeprecateLatestAccessToken(ctx, s.deps.Database, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	s.deps.AccessTokens.Delete(tokenId)
	err = database.RotateUserSecret(ctx, s.deps.Database, user.Id, common.HashSha256(userSecret))
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{
		UserId:     user.Id,
		UserSecret: userSecret,
		StatusCode: 201,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return err
	}
	tokenId, err := database.DeprecateLatestAccessToken(ctx, s.deps.Database, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	s.deps.AccessTokens.Delete(tokenId)
	err = database.DeleteUser(ctx, s.deps.Database, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	s.deps.AuthUsers.Delete(userId)
	return nil
}

func (s *service) GetAccessToken(ctx context.Context) (*common.AccessToken, error) {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	tokenId, err := database.DeprecateLatestAccessToken(ctx, tx, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	s.deps.AccessTokens.Delete(tokenId)
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
	token := &common.AccessToken{
		Id:        accessToken,
		UserId:    user.Id,
		IsLatest:  true,
		InvalidAt: invalidAt,
	}
	s.deps.AccessTokens.Add(tokenId, token)
	return token, nil
}

func (s *service) CreateSecret(ctx context.Context, createArgs *CreateSecretRequest) (*common.Secret, error) {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.deps.SecretsManager.LogAccess(ctx, common.NewAccessLog(user.Id, "CreateSecret", createArgs.Name))
	if err != nil {
		return nil, err
	}
	secretId := common.GenUuid()
	ciphertext, encryptedSecret, err := common.EncryptSecret(secretId, createArgs.Value, common.KEY_SIZE)
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
		StatusCode:  201,
	}
	err = database.CreateSecret(ctx, s.deps.Database, secret)
	if err != nil {
		return nil, err
	}
	secret.Value = createArgs.Value
	secret.CreatedBy = user.Name
	secret.UpdatedBy = user.Name
	return secret, nil
}

func (s *service) GetSecret(ctx context.Context, secretName string) (*common.Secret, error) {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.deps.SecretsManager.LogAccess(ctx, common.NewAccessLog(user.Id, "GetSecret", secretName))
	if err != nil {
		return nil, err
	}

	dbSecret, err := database.GetSecretByName(ctx, s.deps.Database, user, secretName)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := s.deps.SecretsManager.GetSecret(ctx, dbSecret.Id)
	if err != nil {
		return nil, err
	}

	plaintext, err := common.DecryptSecret(dbSecret.Value, encryptedSecret)
	if err != nil {
		return nil, err
	}

	dbSecret.Value = plaintext
	return dbSecret, nil
}

func (s *service) DeleteSecret(ctx context.Context, secretName string) error {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return err
	}
	err = s.deps.SecretsManager.LogAccess(ctx, common.NewAccessLog(user.Id, "DeleteSecret", secretName))
	if err != nil {
		return err
	}

	err = database.DeleteSecret(ctx, s.deps.Database, user.Id, secretName)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GrantPermission(ctx context.Context, req *SecretPermissionRequest) error {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.deps.SecretsManager.LogAccess(ctx, common.NewAccessLog(user.Id, "DeletePermission", req.SecretName))
	if err != nil {
		return err
	}

	secretId, err := database.GetSecretIdWithWriteAccess(ctx, s.deps.Database, user, req.SecretName)
	if err != nil {
		return err
	}

	return database.CreateSecretPermission(ctx, s.deps.Database, user.Id, req.UserId, secretId)
}

func (s *service) RevokePermission(ctx context.Context, req *SecretPermissionRequest) error {
	user, err := common.FetchUserFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.deps.SecretsManager.LogAccess(ctx, common.NewAccessLog(user.Id, "RevokePermission", req.SecretName))
	if err != nil {
		return err
	}

	secretId, err := database.GetSecretIdWithWriteAccess(ctx, s.deps.Database, user, req.SecretName)
	if err != nil {
		return err
	}
	return database.DeleteSecretPermission(ctx, s.deps.Database, user.Id, req.UserId, secretId)
}
