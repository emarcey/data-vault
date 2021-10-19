package handlers

import (
	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/dependencies"
)

func HandleAdminEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointClientAuthenticationWrapper(e, op, deps, true), op, deps), op, deps)
}

func HandleClientEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointClientAuthenticationWrapper(e, op, deps, false), op, deps), op, deps)
}

func HandleTokenEndpoints(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return EndpointLoggingWrapper(EndpointTracingWrapper(EndpointAccessTokenAuthenticationWrapper(e, op, deps), op, deps), op, deps)
}
