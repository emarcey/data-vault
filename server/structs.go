package server

type CreateUserRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateUserResponse struct {
	UserId     string `json:"user_id"`
	UserSecret string `json:"user_secret"`
}

type CreateTableRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteTablePermissionRequest struct {
	UserId  string `json:"user_id"`
	TableId string `json:"table_id"`
}

type CreateTablePermissionRequest struct {
	UserId           string `json:"user_id"`
	TableId          string `json:"table_id"`
	IsDecryptAllowed bool   `json:"is_decrypt_allowed"`
}
