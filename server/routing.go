package server

import (
	"net/http"

	// httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"emarcey/data-vault/dependencies"
	// "emarcey/data-vault/server/handlers"
)

func MakeHttpHandler(s Service, deps *dependencies.Dependencies) http.Handler {
	r := mux.NewRouter()
	// options := []httptransport.ServerOption{
	// 	httptransport.ServerErrorEncoder(handlers.EncodeError),
	// 	httptransport.ServerBefore(handlers.WriteHeadersToContext()),
	// }

	r.Methods("GET").Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(s.Version()))
	})
	return r
}
