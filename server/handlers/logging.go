package handlers

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/emarcey/data-vault/dependencies"
)

// EndpointLoggingWrapper adds logs to every endpoint, exactly like it sounds
func EndpointLoggingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if deps.Env == "local" {
			deps.Logger.Debugf("%s called with %v", op, req)
		}
		resp, err := e(ctx, req)
		if err != nil {
			deps.Logger.Errorf("Endpoint %s returned with error: %v", op, err)
			return nil, err
		}
		if deps.Env == "local" {
			deps.Logger.Debugf("%s returned %v", op, resp)
		}
		return resp, err
	}
}
