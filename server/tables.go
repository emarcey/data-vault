package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func getTableEndpoint(s Service) endpointBuilder {
	op := "GetTable"
	e := func(ctx context.Context, tableIdInterface interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		tableId, ok := tableIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected table ID of type string. Got %T", tableIdInterface)
		}
		return s.GetTable(ctx, user, tableId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "GET",
		path:     "/tables/{id}",
	}
}

func deleteTableEndpoint(s Service) endpointBuilder {
	op := "DeleteTable"
	e := func(ctx context.Context, tableIdInterface interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		tableId, ok := tableIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected table ID of type string. Got %T", tableIdInterface)
		}
		return nil, s.DeleteTable(ctx, user, tableId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "DELETE",
		path:     "/tables/{id}",
	}
}

func decodeCreateTableRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req CreateTableRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("CreateTable", "Could not unmarshal request: %v", string(data))
	}
	fmt.Printf("req: %v\n", req)
	return &req, nil
}

func createTableEndpoint(s Service) endpointBuilder {
	op := "CreateTable"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		req, ok := reqInterface.(*CreateTableRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateTableRequest. Got %T", reqInterface)
		}
		return s.CreateTable(ctx, user, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeCreateTableRequest,
		method:   "POST",
		path:     "/tables",
	}
}
