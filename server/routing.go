package server

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"emarcey/data-vault/dependencies"
	"emarcey/data-vault/server/handlers"
)

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
	r.Methods("GET").Path("/users").Handler(httptransport.NewServer(
		handlers.HandleAdminEndpoints(listUsersEndpoint(s), "ListUsers", deps),
		decodeListUsersRequest,
		encodeResponse,
		options...,
	))
	return r
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// Don't overwrite a header (i.e. called from encodeTextResponse)
	if v := w.Header().Get("Content-Type"); v == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Only write json body if we're setting response as json
		return json.NewEncoder(w).Encode(response)
	}
	return nil
}
