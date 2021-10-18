package handlers

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/dependencies"
)

func EndpointTracingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()
		resp, err := e(tracer.Context(), request)
		return resp, err
	}
}
