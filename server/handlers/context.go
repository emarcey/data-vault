package handlers

import (
	"context"
	"net/http"

	"emarcey/data-vault/common"
	httptransport "github.com/go-kit/kit/transport/http"
)

func WriteHeadersToContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, common.HeadersContextKey, r.Header)
	}
}
