package handlers

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/emarcey/data-vault/dependencies"
)

// EndpointLoggingWrapper adds logs to every endpoint, exactly like it sounds
func EndpointLoggingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := e(ctx, request)
		if err != nil {
			deps.Logger.Errorf("Endpoint %s returned with error: %v", op, err)
			return nil, err
		}
		return resp, err
	}
}
