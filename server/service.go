package server

import (
	"context"

	"emarcey/data-vault/common"
	"emarcey/data-vault/database"
	"emarcey/data-vault/dependencies"
)

type Service interface {
	Version() string
	ListUsers(ctx context.Context) ([]*common.User, error)
	GetUser(ctx context.Context, id string) (*common.User, error)
	CreateUser(ctx context.Context, name string, userType string) (*common.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetAccessToken(ctx context.Context, id string, clientId string) (*common.AccessToken, error)
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

func (s *service) GetUser(ctx context.Context, id string) (*common.User, error) {
	return nil, nil
}

func (s *service) CreateUser(ctx context.Context, name string, userType string) (*common.User, error) {
	return nil, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (s *service) GetAccessToken(ctx context.Context, id string, clientId string) (*common.AccessToken, error) {
	return nil, nil
}
