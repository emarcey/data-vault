package server

import (
	"context"
	// "encoding/json"
	// "io/ioutil"
	// "net/http"

	// "github.com/gorilla/mux"

	"emarcey/data-vault/common"
)

func listTablesEndpoint(s Service) endpointBuilder {
	e := func(ctx context.Context, _ interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		return s.ListTables(ctx, user)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  noOpDecodeRequest,
		method:   "GET",
		path:     "/tables",
	}
}
