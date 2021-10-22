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

	// tables
	ListTables(ctx context.Context, user *common.User) ([]*common.Table, error)
	GetTable(ctx context.Context, user *common.User, tableId string) (*common.Table, error)
	DeleteTable(ctx context.Context, user *common.User, tableId string) error
	CreateTable(ctx context.Context, user *common.User, req *CreateTableRequest) (*common.Table, error)

	// table permissions
	ListTablePermissions(ctx context.Context, userId string) ([]*common.TablePermission, error)
	ListTablePermissionsForTable(ctx context.Context, tableId string) ([]*common.TablePermission, error)
	DeleteTablePermission(ctx context.Context, adminUser *common.User, req *DeleteTablePermissionRequest) error
	CreateTablePermission(ctx context.Context, adminUser *common.User, req *CreateTablePermissionRequest) (*common.TablePermission, error)
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

func (s *service) ListTables(ctx context.Context, user *common.User) ([]*common.Table, error) {
	return database.ListTables(ctx, s.deps.Database, user)
}

func (s *service) GetTable(ctx context.Context, user *common.User, tableId string) (*common.Table, error) {
	return database.GetTableById(ctx, s.deps.Database, user, tableId)
}

func (s *service) DeleteTable(ctx context.Context, user *common.User, tableId string) error {
	return database.DeleteTable(ctx, s.deps.Database, user, tableId)
}

func (s *service) CreateTable(ctx context.Context, user *common.User, req *CreateTableRequest) (*common.Table, error) {
	tx, err := s.deps.Database.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	return database.CreateTable(ctx, tx, user, req.Name, req.Description)
}

func (s *service) ListTablePermissions(ctx context.Context, userId string) ([]*common.TablePermission, error) {
	return database.ListTablePermissions(ctx, s.deps.Database, userId)
}

func (s *service) ListTablePermissionsForTable(ctx context.Context, tableId string) ([]*common.TablePermission, error) {
	return database.ListTablePermissionsForTable(ctx, s.deps.Database, tableId)
}

func (s *service) DeleteTablePermission(ctx context.Context, adminUser *common.User, req *DeleteTablePermissionRequest) error {
	return database.DeleteTablePermission(ctx, s.deps.Database, adminUser, req.UserId, req.TableId)
}

func (s *service) CreateTablePermission(ctx context.Context, adminUser *common.User, req *CreateTablePermissionRequest) (*common.TablePermission, error) {
	return database.CreateTablePermission(ctx, s.deps.Database, adminUser, req.UserId, req.TableId, req.IsDecryptAllowed)
}
