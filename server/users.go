package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"emarcey/data-vault/common"
)

func listUsersEndpoint(s Service) endpointBuilder {
	e := func(ctx context.Context, _ interface{}) (interface{}, error) {
		return s.ListUsers(ctx)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  noOpDecodeRequest,
		method:   "GET",
		path:     "/users",
	}
}

func getUserEndpoint(s Service) endpointBuilder {
	op := "GetUser"
	e := func(ctx context.Context, userIdInterface interface{}) (interface{}, error) {
		userId, ok := userIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user ID of type string. Got %T", userIdInterface)
		}
		return s.GetUser(ctx, userId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "GET",
		path:     "/users/{id}",
	}
}

func deleteUserEndpoint(s Service) endpointBuilder {
	op := "DeleteUser"
	e := func(ctx context.Context, userIdInterface interface{}) (interface{}, error) {
		userId, ok := userIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user ID of type string. Got %T", userIdInterface)
		}
		return nil, s.DeleteUser(ctx, userId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "DELETE",
		path:     "/users/{id}",
	}
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req CreateUserRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("CreateUser", "Could not unmarshal request: %v", string(data))
	}
	return &req, nil
}

func createUserEndpoint(s Service) endpointBuilder {
	op := "CreateUser"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*CreateUserRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateUserRequest. Got %T", reqInterface)
		}
		return s.CreateUser(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeCreateUserRequest,
		method:   "POST",
		path:     "/users",
	}
}

func getAccessTokenEndpoint(s Service) endpointBuilder {
	e := func(ctx context.Context, _ interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		return s.GetAccessToken(ctx, user)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  noOpDecodeRequest,
		method:   "GET",
		path:     "/access_token",
	}
}
