package common

import (
	"time"
)

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	Type       string `json:"type"`
	SecretHash string `json:"-"`
}

type AccessToken struct {
	Id        string    `json:"id"`
	UserId    string    `json:"client_id"`
	InvalidAt time.Time `json:"invalid_at"`
	IsLatest  bool      `json:"is_latest"`
}

type Secret struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}

type EncryptedSecret struct {
	Id  string `json:"_id" bson:"_id"`
	Key string `json:"key" bson:"key"`
	Iv  string `json:"iv" bson:"iv"`
}
