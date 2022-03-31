package handlers

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/common"
	"emarcey/data-vault/common/tracer"
	"emarcey/data-vault/dependencies"
)

func authenticateClient(ctx context.Context, op string, tracer tracer.Tracer, authUsers *dependencies.UserCache, checkAdmin bool) (*common.User, error) {
	userId, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_ID)
	if err != nil {
		return nil, err
	}
	tracer.AddBreadcrumb(map[string]interface{}{"userId": userId})

	userSecretRaw, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_CLIENT_SECRET)
	if err != nil {
		return nil, err
	}

	userSecret := common.HashSha256(userSecretRaw)

	user := authUsers.Get(userId)
	if user == nil {
		return nil, fmt.Errorf("User not found for userId %s", userId)
	}

	if user.SecretHash != userSecret {
		return nil, fmt.Errorf("Invalid secret for userId %s", userId)
	}

	if checkAdmin && user.Type != "admin" {
		return nil, fmt.Errorf("User %s is not an admin.", user.Id)
	}
	return user, nil
}

// EndpointClientAuthenticationWrapper validates request authentication by client id/secret, with admin check optional
func EndpointClientAuthenticationWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies, checkAdmin bool) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()

		user, err := authenticateClient(ctx, op, tracer, deps.AuthUsers, checkAdmin)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		newCtx := common.InjectUserIntoContext(tracer.Context(), user)
		return e(newCtx, request)
	}
}

func authenticateAccessToken(
	ctx context.Context,
	op string,
	tracer tracer.Tracer,
	authUsers *dependencies.UserCache,
	accessTokens *dependencies.AccessTokenCache,
	checkAdmin bool,
) (*common.User, error) {
	authTokenRaw, err := common.FetchStringFromContextHeaders(ctx, common.HEADER_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	authToken := common.HashSha256(authTokenRaw)
	tracer.AddBreadcrumb(map[string]interface{}{"authToken": authToken})

	accessToken := accessTokens.Get(authToken)
	if accessToken == nil {
		return nil, fmt.Errorf("Auth Token not found")
	}

	tracer.AddBreadcrumb(map[string]interface{}{"userId": accessToken.UserId})

	user := authUsers.Get(accessToken.UserId)
	if user == nil {
		return nil, fmt.Errorf("User not found for accessToken %s", authToken)
	}

	if checkAdmin && user.Type != "admin" {
		return nil, fmt.Errorf("User %s is not an admin.", user.Id)
	}
	return user, nil
}

// EndpointAccessTokenAuthenticationWrapper validates request authentication by access token
func EndpointAccessTokenAuthenticationWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies, checkAdmin bool) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()

		user, err := authenticateAccessToken(ctx, op, tracer, deps.AuthUsers, deps.AccessTokens, checkAdmin)
		if err != nil {
			tracer.CaptureException(err)
			deps.Logger.Errorf("Error authenticating %s: %v", op, err)
			return nil, common.NewAuthorizationError()
		}
		newCtx := common.InjectUserIntoContext(tracer.Context(), user)
		return e(newCtx, request)
	}
}
