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

type Table struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}

type Column struct {
	TableId    string `json:"table_id"`
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

type TablePermission struct {
	UserId           string `json:"client_id"`
	TableId          string `json:"table_id"`
	TableName        string `json:"table_name"`
	IsDecryptAllowed bool   `json:"is_decrypt_allowed"`
	CreatedBy        string `json:"created_by"`
	UpdatedBy        string `json:"updated_by"`
}
