package handlers

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/emarcey/data-vault/dependencies"
)

// EndpointTracingWrapper adds a tracing context to every endpoint, exactly like it sounds
func EndpointTracingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracer := deps.Tracer(ctx, op)
		defer tracer.Close()
		resp, err := e(tracer.Context(), request)
		return resp, err
	}
}
