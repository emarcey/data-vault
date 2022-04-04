package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emarcey/data-vault/common"
)

var decodeSecretPermissionUrl = decodeRequestUrlName("SecretPermission")

func decodeSecretPermissionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	op := "SecretPermission"
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req SecretPermissionRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError(op, "Could not unmarshal request: %v", string(data))
	}
	secretName, err := decodeSecretPermissionUrl(ctx, r)
	if err != nil {
		return nil, err
	}
	req.SecretName = secretName.(string)
	return &req, nil
}

func createSecretPermissionEndpoint(s Service) endpointBuilder {
	op := "CreateSecretPermission"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*SecretPermissionRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *SecretPermissionRequest. Got %T", reqInterface)
		}
		err := s.GrantPermission(ctx, req)
		if err != nil {
			return nil, err
		}
		return NewStatusResponse(), nil
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeSecretPermissionRequest,
		method:   HTTP_POST,
		path:     "/secrets/{name}/permissions",
	}
}

func deleteSecretPermissionEndpoint(s Service) endpointBuilder {
	op := "RevokeSecretPermission"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*SecretPermissionRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *SecretPermissionRequest. Got %T", reqInterface)
		}
		return nil, s.RevokePermission(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeSecretPermissionRequest,
		method:   HTTP_DELETE,
		path:     "/secrets/{name}/permissions",
	}
}
