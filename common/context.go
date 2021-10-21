package common

import (
	"context"
	"net/http"
)

var HeadersContextKey struct{}
var UserContextKey struct{}

func FetchStringFromContextHeaders(ctx context.Context, key string) (string, error) {
	op := "FetchStringFromContextHeaders"
	headersInterface := ctx.Value(HeadersContextKey)
	if headersInterface == nil {
		return "", NewInvalidParamsError(op, "No headers in context.")
	}
	headers, ok := headersInterface.(http.Header)
	if !ok {
		return "", NewInvalidParamsError(op, "Expected headers of type http.Header. Got %T", headersInterface)
	}

	val, ok := headers[key]
	if !ok {
		return "", NewInvalidParamsError(op, "Key %s not in header", key)
	}
	if len(val) != 1 {
		return "", NewInvalidParamsError(op, "Expected 1 value for key, %s. Got %d", key, len(val))
	}

	return val[0], nil
}

func InjectUserIntoContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

func FetchUserFromContext(ctx context.Context) (*User, error) {
	op := "FetchUserFromContext"
	userInterface := ctx.Value(UserContextKey)
	if userInterface == nil {
		return nil, NewInvalidParamsError(op, "No user in context.")
	}
	user, ok := userInterface.(*User)
	if !ok {
		return nil, NewInvalidParamsError(op, "Expected user of type *common.User. Got %T", userInterface)
	}
	if user == nil {
		return nil, NewInvalidParamsError(op, "User in context was nil")
	}

	return user, nil
}
