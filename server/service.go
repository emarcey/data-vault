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
	ListUsers(ctx context.Context) ([]*common.User, error)
	GetUser(ctx context.Context, id string) (*common.User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	GetAccessToken(ctx context.Context, userId string) (*common.AccessToken, error)
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

type CreateUserRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateUserResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	clientId := common.GenUuid()
	clientSecret := common.GenUuid()
	err := database.CreateUser(ctx, s.deps.Database, clientId, req.Name, req.Type, common.HashSha256(clientSecret))
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, userId string) error {
	err := database.DeprecateLatestAccessToken(ctx, s.deps.Database, userId)
	if err != nil {
		return err
	}
	return database.DeleteUser(ctx, s.deps.Database, userId)
}

func (s *service) GetAccessToken(ctx context.Context, userId string) (*common.AccessToken, error) {
	err := database.DeprecateLatestAccessToken(ctx, s.deps.Database, userId)
	if err != nil {
		return nil, err
	}
	accessToken := common.GenUuid()
	invalidAt := time.Now().Add(time.Duration(s.deps.ServerConfigs.AccessTokenHours) * time.Hour)
	err = database.CreateAccessToken(ctx, s.deps.Database, userId, common.HashSha256(accessToken), invalidAt)
	if err != nil {
		return nil, err
	}
	return &common.AccessToken{
		Id:        accessToken,
		UserId:    userId,
		IsLatest:  true,
		InvalidAt: invalidAt,
	}, nil
}
