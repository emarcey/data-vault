package server

import (
	"emarcey/data-vault/common"
)

type Service interface {
	ListUsers() ([]*common.User, error)
	GetUser(id string) (*common.User, error)
	CreateUser(name string, userType string) (*common.User, error)
	DeleteUser(id string) error
	GetAccessToken(id string, clientId string) (*common.AccessToken, error)
}

func (s *Service) ListUsers() ([]*common.User, error) {

}
func (s *Service) GetUser(id string) (*common.User, error) {

}
func (s *Service) CreateUser(name string, userType string) (*common.User, error) {

}
func (s *Service) DeleteUser(id string) error {

}
func (s *Service) GetAccessToken(id string, clientId string) (*common.AccessToken, error) {

}
