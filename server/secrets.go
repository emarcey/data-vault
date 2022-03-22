package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"emarcey/data-vault/common"
)

func decodeCreateSecretRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req CreateSecretRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, common.NewInvalidParamsError("CreateSecret", "Could not unmarshal request: %v", string(data))
	}
	return &req, nil
}

func createSecretEndpoint(s Service) endpointBuilder {
	op := "CreateSecret"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*CreateSecretRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *CreateSecretRequest. Got %T", reqInterface)
		}
		return s.CreateSecret(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeCreateSecretRequest,
		method:   "POST",
		path:     "/secrets",
	}
}
