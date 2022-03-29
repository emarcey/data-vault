package server

type CreateUserRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateUserGroupRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	UserId     string `json:"user_id"`
	UserSecret string `json:"user_secret"`
	StatusCode int    `json:"-"`
}

func (c *CreateUserResponse) GetStatusCode() int {
	if c.StatusCode == 0 {
		return 200
	}
	return c.StatusCode
}

type CreateSecretRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type SecretPermissionRequest struct {
	SecretName  string `json:"-"`
	UserId      string `json:"user_id"`
	UserGroupId string `json:"user_group_id"`
}

type UserGroupMemberRequest struct {
	UserGroupId string `json:"-"`
	UserId      string `json:"user_id"`
}

type SimpleCreateResponse struct {
	StatusCode int `json:"-"`
}

func NewSimpleCreateResponse() *SimpleCreateResponse {
	return &SimpleCreateResponse{StatusCode: 201}
}

func (c *SimpleCreateResponse) GetStatusCode() int {
	if c.StatusCode == 0 {
		return 201
	}
	return c.StatusCode
}
