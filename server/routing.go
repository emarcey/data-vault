package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"emarcey/data-vault/dependencies"
)

func MakeHttpHandler(s Service, deps *dependencies.Dependencies) http.Handler {

}
