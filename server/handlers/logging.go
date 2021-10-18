package handlers

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/dependencies"
)

func EndpointLoggingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		deps.Logger.Infof("Endpoint %s called with request: %v", op, request)
		defer tracer.Close()
		resp, err := e(tracer.Context(), request)
		if err != nil {
			deps.Logger.Errorf("Endpoint %s returned with error: %v", op, err)
			return nil, err
		}
		deps.Logger.Infof("Endpoint %s returned: %v", op, resp)
		return resp, err
	}
}
