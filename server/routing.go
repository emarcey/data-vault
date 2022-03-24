package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"emarcey/data-vault/common"
	"emarcey/data-vault/dependencies"
	"emarcey/data-vault/server/handlers"
)

type endpointBuilder struct {
	endpoint endpoint.Endpoint
	decoder  httptransport.DecodeRequestFunc
	method   string
	path     string
}

func makeMethods(r *mux.Router, deps *dependencies.Dependencies, handler handlers.EndpointHandler, endpoints []endpointBuilder, encoder httptransport.EncodeResponseFunc, options ...httptransport.ServerOption) {
	for _, endpoint := range endpoints {
		r.Methods(endpoint.method).Path(endpoint.path).Handler(httptransport.NewServer(
			handler(endpoint.endpoint, fmt.Sprintf("%s %s", endpoint.method, endpoint.path), deps),
			endpoint.decoder,
			encoder,
			options...,
		))
	}
}

func MakeHttpHandler(s Service, deps *dependencies.Dependencies) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(handlers.EncodeError),
		httptransport.ServerBefore(handlers.WriteHeadersToContext()),
	}

	r.Methods("GET").Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(s.Version()))
	})
	adminEndpoints := []endpointBuilder{
		getUserEndpoint(s),
		deleteUserEndpoint(s),
		createUserEndpoint(s),
		deleteSecretEndpoint(s),
	}
	makeMethods(r, deps, handlers.HandleAdminEndpoints, adminEndpoints, encodeResponse, options...)

	clientEndpoints := []endpointBuilder{
		getAccessTokenEndpoint(s),
	}
	makeMethods(r, deps, handlers.HandleClientEndpoints, clientEndpoints, encodeResponse, options...)

	accessTokenEndpoints := []endpointBuilder{
		listUsersEndpoint(s),
		createSecretEndpoint(s),
		getSecretEndpoint(s),
		createSecretPermissionEndpoint(s),
		deleteSecretPermissionEndpoint(s),
	}
	makeMethods(r, deps, handlers.HandleTokenEndpoints, accessTokenEndpoints, encodeResponse, options...)
	return r
}

func noOpDecodeRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeRequestUrlId(op string) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Id not found")
		}
		return id, nil
	}
}

func decodeRequestUrlName(op string) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		id, ok := vars["name"]
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Name not found")
		}
		return id, nil
	}
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		w.WriteHeader(204)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	responser, ok := response.(common.Responser)
	if ok {
		w.WriteHeader(responser.GetStatusCode())
	}
	return json.NewEncoder(w).Encode(response)
}
