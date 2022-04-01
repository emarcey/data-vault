package server

type PaginationRequest struct {
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

type ListUsersInGroupRequest struct {
	UserGroupId string
	PageSize    int `json:"page_size"`
	Offset      int `json:"offset"`
}

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

type StatusResponse struct {
	StatusCode int `json:"-"`
}

func NewStatusResponse() *StatusResponse {
	return &StatusResponse{StatusCode: 201}
}

func (c *StatusResponse) GetStatusCode() int {
	if c.StatusCode == 0 {
		return 201
	}
	return c.StatusCode
}
