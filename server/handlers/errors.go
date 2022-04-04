package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emarcey/data-vault/common"
)

// EncodeError JSON encodes the supplied error
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		err = fmt.Errorf("EncodeError called with nil error")
	}
	errorCode := 500
	newErr, ok := err.(common.ErrorWithCode)
	if ok {
		errorCode = newErr.Code()
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
