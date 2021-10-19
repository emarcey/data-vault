package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func decodeListUsersRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func listUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return s.ListUsers(ctx)
	}
}
