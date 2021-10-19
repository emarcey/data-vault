package handlers

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/common"
	"emarcey/data-vault/dependencies"
)

func EndpointClientAuthenticationWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies, checkAdmin bool) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()
		clientId, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_ID)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		clientSecretRaw, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_SECRET)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		clientSecret := common.HashSha256(clientSecretRaw)

		user, ok := deps.AuthUsers[fmt.Sprintf(`%s_%s`, clientId, clientSecret)]
		if !ok {
			internalError := fmt.Errorf("User not found for clientId %s", clientId)
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}

		if checkAdmin && user.Type != "admin" {
			internalError := fmt.Errorf("User %s is not an admin.", user.Id)
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}
		return e(tracer.Context(), request)
	}
}

func EndpointAccessTokenAuthenticationWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()

		authTokenRaw, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_ACCESS_TOKEN)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		authToken := common.HashSha256(authTokenRaw)
		_, ok := deps.AccessTokens[authToken]
		if !ok {
			internalError := fmt.Errorf("Auth Token not found")
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}
		return e(tracer.Context(), request)
	}
}
