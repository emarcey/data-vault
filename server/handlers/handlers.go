package handlers

import (
	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/dependencies"
)

type EndpointHandler func(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint

// HandleAdminEndpoints -- wrapper to add logging/tracing/auth for user_id/secret admin endpoints
func HandleAdminEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointClientAuthenticationWrapper(e, op, deps, true), op, deps), op, deps)
}

// HandleClientEndpoints -- wrapper to add logging/tracing/auth for user_id/secret endpoints
func HandleClientEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointClientAuthenticationWrapper(e, op, deps, false), op, deps), op, deps)
}

// HandleTokenEndpoints -- wrapper to add logging/tracing/auth for access_token endpoints
func HandleTokenEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointAccessTokenAuthenticationWrapper(e, op, deps), op, deps), op, deps)
}
