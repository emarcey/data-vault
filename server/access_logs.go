package server

import (
	"context"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/emarcey/data-vault/common"
)

func decodeAccessLogsRequest(op string) httptransport.DecodeRequestFunc {
	userIdDecoder := decodeRequestUrlId(op)
	paginationDecoder := decodePaginationRequest(op)
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		userIdInterface, err := userIdDecoder(ctx, r)
		if err != nil {
			return nil, err
		}
		userId, ok := userIdInterface.(string)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected id of type string Got %T", userIdInterface)
		}

		paginationInterface, err := paginationDecoder(ctx, r)
		if err != nil {
			return nil, err
		}
		pagination := paginationInterface.(*PaginationRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected pagination of type string Got %T", paginationInterface)
		}

		urlParams := r.URL.Query()
		startDate, err := parseDateUrlParam(op, urlParams, "startDate", common.DEFAULT_START_TIME)
		if err != nil {
			return nil, err
		}

		endDate, err := parseDateUrlParam(op, urlParams, "endDate", time.Now())
		if err != nil {
			return nil, err
		}

		return &common.ListAccessLogsRequest{
			UserId:    userId,
			PageSize:  pagination.PageSize,
			Offset:    pagination.Offset,
			StartDate: startDate,
			EndDate:   endDate,
		}, nil
	}
}

func listAccessLogsEndpoint(s Service) endpointBuilder {
	op := "ListAccessLogs"
	e := func(ctx context.Context, reqInterface interface{}) (interface{}, error) {
		req, ok := reqInterface.(*common.ListAccessLogsRequest)
		if !ok {
			return nil, common.NewInvalidParamsError(op, "Expected request of type *common.ListAccessLogsRequest Got %T", reqInterface)
		}
		return s.ListAccessLogs(ctx, req)
	}
	return endpointBuilder{
		endpoint: e,
		decoder:  decodeAccessLogsRequest(op),
		method:   HTTP_GET,
		path:     "/users/{id}/access-logs",
	}
}
