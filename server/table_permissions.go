package server

import (
	"context"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"

	// "github.com/gorilla/mux"

	"emarcey/data-vault/common"
)

func listTablePermissionsEndpoint(s Service) endpointBuilder {
	e := func(ctx context.Context, _ interface{}) (interface{}, error) {
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		return s.ListTablePermissions(ctx, user.Id)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  noOpDecodeRequest,
		method:   "GET",
		path:     "/table-permissions",
	}
}

func listTablePermissionsForUserEndpoint(s Service) endpointBuilder {
	op := "ListTablePermissionsForUser"
	e := func(ctx context.Context, userIdInterface interface{}) (interface{}, error) {
		userId, ok := userIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected user ID of type string. Got %T", userIdInterface)
		}
		return s.ListTablePermissions(ctx, userId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "GET",
		path:     "/users/{id}/table-permissions",
	}
}

func listTablePermissionsForTableEndpoint(s Service) endpointBuilder {
	op := "ListTablePermissionsForTable"
	e := func(ctx context.Context, tableIdInterface interface{}) (interface{}, error) {
		tableId, ok := tableIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected table ID of type string. Got %T", tableIdInterface)
		}
		return s.ListTablePermissionsForTable(ctx, tableId)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeRequestUrlId(op),
		method:   "GET",
		path:     "/tables/{id}/table-permissions",
	}
}

func decodeDeleteTablePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req DeleteTablePermissionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("DeleteTablePermission", "Could not unmarshal request: %v", string(data))
	}
	return &req, nil
}

func deleteTablePermissionEndpoint(s Service) endpointBuilder {
	op := "DeleteTablePermission"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*DeleteTablePermissionRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *DeleteTablePermissionRequest. Got %T", reqInterface)
		}
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		return nil, s.DeleteTablePermission(ctx, user, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeDeleteTablePermissionRequest,
		method:   "DELETE",
		path:     "/table-permissions",
	}
}

func decodeCreateTablePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req CreateTablePermissionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("CreateTablePermission", "Could not unmarshal request: %v", string(data))
	}
	return &req, nil
}

func createTablePermissionEndpoint(s Service) endpointBuilder {
	op := "CreateTablePermission"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*CreateTablePermissionRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateTablePermissionRequest. Got %T", reqInterface)
		}
		user, err := common.FetchUserFromContext(ctx)
		if err != nil {
			return nil, err
		}
		return s.CreateTablePermission(ctx, user, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeCreateTablePermissionRequest,
		method:   "POST",
		path:     "/table-permissions",
	}
}

// func decodeCreateTableRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	data, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var req CreateTableRequest
// 	err = json.Unmarshal(data, &req)
// 	if err != nil {
// 		return nil, common.NewInvalidParamsError("CreateTable", "Could not unmarshal request: %v", string(data))
// 	}
// 	fmt.Printf("req: %v\n", req)
// 	return &req, nil
// }

// func createTableEndpoint(s Service) endpointBuilder {
// 	op := "CreateTable"
// 	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
// 		user, err := common.FetchUserFromContext(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		req, ok := reqInterface.(*CreateTableRequest)
// 		if !ok {
// 			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateTableRequest. Got %T", reqInterface)
// 		}
// 		return s.CreateTable(ctx, user, req)
// 	}
// 	return endpointBuilder{
// 		endpoint: e,
// 		decoder:  decodeCreateTableRequest,
// 		method:   "POST",
// 		path:     "/tables",
// 	}
// }
