package handlers

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/common"
	"emarcey/data-vault/dependencies"
)

// EndpointClientAuthenticationWrapper validates request authentication by client id/secret, with admin check optional
func EndpointClientAuthenticationWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies, checkAdmin bool) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()

		userId, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_ID)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}

		userSecretRaw, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_SECRET)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		userSecret := common.HashSha256(userSecretRaw)

		user, ok := deps.AuthUsers[userId]
		if !ok {
			internalError := fmt.Errorf("User not found for userId %s", userId)
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}

		if user.SecretHash != userSecret {
			internalError := fmt.Errorf("Invalid secret for userId %s", userId)
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
		newCtx := common.InjectUserIntoContext(tracer.Context(), user)
		return e(newCtx, request)
	}
}

// EndpointAccessTokenAuthenticationWrapper validates request authentication by access token
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

		accessToken, ok := deps.AccessTokens[authToken]
		if !ok {
			internalError := fmt.Errorf("Auth Token not found")
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}

		user, ok := deps.AuthUsers[accessToken.UserId]
		if !ok {
			internalError := fmt.Errorf("User not found for accessToken %s", authToken)
			tracer.CaptureException(internalError)
			deps.Logger.Errorf("Error authenticating %s: %v", op, internalError)
			return nil, common.NewAuthorizationError()
		}
		newCtx := common.InjectUserIntoContext(tracer.Context(), user)
		return e(newCtx, request)
	}
}
