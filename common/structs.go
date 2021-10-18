package common

import (
	"time"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Type     string `json:"type"`
}

type AccessToken struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	InvalidAt time.Time `json:"invalid_at"`
	IsLatest  bool      `json:"is_latest"`
}

type TablePermissions struct {
	UserId           string `json:"user_id"`
	TableName        string `json:"table_name"`
	IsDecryptAllowed bool   `json:"is_decrypt_allowed"`
	IsActive         bool   `json:"is_active"`
}
