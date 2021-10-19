package server

import (
	"emarcey/data-vault/common"
	"emarcey/data-vault/dependencies"
)

type Service interface {
	Version() string
	ListUsers() ([]*common.User, error)
	GetUser(id string) (*common.User, error)
	CreateUser(name string, userType string) (*common.User, error)
	DeleteUser(id string) error
	GetAccessToken(id string, clientId string) (*common.AccessToken, error)
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

func (s *service) ListUsers() ([]*common.User, error) {
	return nil, nil
}

func (s *service) GetUser(id string) (*common.User, error) {
	return nil, nil
}

func (s *service) CreateUser(name string, userType string) (*common.User, error) {
	return nil, nil
}

func (s *service) DeleteUser(id string) error {
	return nil
}

func (s *service) GetAccessToken(id string, clientId string) (*common.AccessToken, error) {
	return nil, nil
}
