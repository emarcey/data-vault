package server

type CreateUserRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateUserResponse struct {
	UserId     string `json:"user_id"`
	UserSecret string `json:"user_secret"`
}

type CreateSecretRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
