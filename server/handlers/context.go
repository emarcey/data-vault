package handlers

import (
	"context"
	"net/http"

	"github.com/emarcey/data-vault/common"
	httptransport "github.com/go-kit/kit/transport/http"
)

// WriteHeadersToContext populates the context with values from the request header
func WriteHeadersToContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return common.InjectHeaderIntoContext(ctx, r)
	}
}
