package handlers

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"

	"emarcey/data-vault/dependencies"
)

// EndpointLoggingWrapper adds logs to every endpoint, exactly like it sounds
func EndpointLoggingWrapper(e endpoint.Endpoint, op string, deps *dependencies.Dependencies) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		deps.Logger.Debugf("Endpoint %s called with request: %v", op, request)
		resp, err := e(ctx, request)
		if err != nil {
			deps.Logger.Errorf("Endpoint %s returned with error: %v", op, err)
			return nil, err
		}
		var logResp interface{}
		marshalledRespBytes, err := json.Marshal(resp)
		if err != nil {
			logResp = resp
		} else {
			logResp = string(marshalledRespBytes)
		}
		deps.Logger.Debugf("Endpoint %s returned: %v", op, logResp)
		return resp, err
	}
}
