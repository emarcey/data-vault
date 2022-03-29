package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"emarcey/data-vault/common"
)

func listUserGroupsEndpoint(s Service) endpointBuilder {
	e := func(ctx context.Context, _ interface{}) (interface{}, error) {
		return s.ListUserGroups(ctx)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  noOpDecodeRequest,
		method:   HTTP_GET,
		path:     "/user-groups",
	}
}

func getUserGroupEndpoint(s Service) endpointBuilder {
	op := "GetUserGroup"
	e := func(ctx context.Context, userGroupIdInterface interface{}) (interface{}, error) {
		userGroupId, ok := userGroupIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user group ID of type string. Got %T", userGroupIdInterface)
		}
		return s.GetUserGroup(ctx, userGroupId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   HTTP_GET,
		path:     "/user-groups/{id}",
	}
}

func deleteUserGroupEndpoint(s Service) endpointBuilder {
	op := "DeleteUserGroup"
	e := func(ctx context.Context, userGroupIdInterface interface{}) (interface{}, error) {
		userGroupId, ok := userGroupIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user group ID of type string. Got %T", userGroupIdInterface)
		}
		return s.GetUserGroup(ctx, userGroupId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   HTTP_DELETE,
		path:     "/user-groups/{id}",
	}
}

func listUsersInGroupEndpoint(s Service) endpointBuilder {
	op := "ListUsersInGroup"
	e := func(ctx context.Context, userGroupIdInterface interface{}) (interface{}, error) {
		userGroupId, ok := userGroupIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user group ID of type string. Got %T", userGroupIdInterface)
		}
		return s.ListUsersInGroup(ctx, userGroupId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   HTTP_GET,
		path:     "/user-groups/{id}/users",
	}
}

func decodeCreateUserGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req CreateUserGroupRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("CreateUserGroup", "Could not unmarshal request: %v", string(data))
	}
	return &req, nil
}

func createUserGroupEndpoint(s Service) endpointBuilder {
	op := "CreateUserGroup"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*CreateUserGroupRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateUserRequest. Got %T", reqInterface)
		}
		return s.CreateUserGroup(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeCreateUserGroupRequest,
		method:   HTTP_POST,
		path:     "/user-groups",
	}
}

var decodeUserGroupMemberRequestId = decodeRequestUrlId("UserGroupMemberRequest")

func decodeUserGroupMemberRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	op := "UserGroupMemberRequest"
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req UserGroupMemberRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError(op, "Could not unmarshal request: %v", string(data))
	}
	userGroupId, err := decodeUserGroupMemberRequestId(ctx, r)
	if err != nil {
		return nil, err
	}
	req.UserGroupId = userGroupId.(string)
	return &req, nil
}

func addUserToGroupEndpoint(s Service) endpointBuilder {
	op := "AddUserToGroup"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*UserGroupMemberRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *UserGroupMemberRequest. Got %T", reqInterface)
		}
		err := s.AddUserToGroup(ctx, req)
		if err != nil {
			return nil, err
		}
		return NewSimpleCreateResponse(), nil
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeUserGroupMemberRequest,
		method:   HTTP_POST,
		path:     "/user-groups/{id}/users",
	}
}

func removeUserFromGroupEndpoint(s Service) endpointBuilder {
	op := "RemoveUserFromGroup"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*UserGroupMemberRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *UserGroupMemberRequest. Got %T", reqInterface)
		}
		return nil, s.RemoveUserFromGroup(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeUserGroupMemberRequest,
		method:   HTTP_DELETE,
		path:     "/user-groups/{id}/users",
	}
}