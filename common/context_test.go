package common

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchStringFromContextHeadersErrors(t *testing.T) {
	var tests = []struct {
		testName string
		ctx      context.Context
		key      string
	}{
		{
			testName: "background context",
			ctx:      context.Background(),
			key:      "mykey",
		},
		{
			testName: "not header",
			ctx:      context.WithValue(context.Background(), HeadersContextKey, "zoop"),
			key:      "mykey",
		},
		{
			testName: "key not found",
			ctx:      context.WithValue(context.Background(), HeadersContextKey, http.Header{"zoop": []string{"zoop"}}),
			key:      "mykey",
		},
		{
			testName: "empty key",
			ctx:      context.WithValue(context.Background(), HeadersContextKey, http.Header{"mykey": []string{}}),
			key:      "mykey",
		},
		{
			testName: "too many values",
			ctx:      context.WithValue(context.Background(), HeadersContextKey, http.Header{"mykey": []string{"1", "2", "3"}}),
			key:      "mykey",
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("FetchStringFromContextHeaders - Errors - %v", given.testName), func(t *testing.T) {
			result, err := FetchStringFromContextHeaders(given.ctx, given.key)
			require.NotNil(t, err, "no error in FetchStringFromContextHeaders: %v", err)
			require.Empty(t, result, "Expected empty result, got: %v", result)
		})
	}
}

func TestFetchStringFromContextHeadersSuccesses(t *testing.T) {
	var tests = []struct {
		testName string
		ctx      context.Context
		key      string
		expected string
	}{
		{
			testName: "too many values",
			ctx:      context.WithValue(context.Background(), HeadersContextKey, http.Header{"mykey": []string{"myvalue"}}),
			key:      "mykey",
			expected: "myvalue",
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("FetchStringFromContextHeaders - Successes - %v", given.testName), func(t *testing.T) {
			result, err := FetchStringFromContextHeaders(given.ctx, given.key)
			require.Nil(t, err, "no error in FetchStringFromContextHeaders: %v", err)
			require.Equal(t, result, given.expected, "Result %v did not equal expected %v", result, given.expected)
		})
	}
}

func TestFetchUserFromContextErrors(t *testing.T) {
	var nilUser *User

	var tests = []struct {
		testName string
		ctx      context.Context
	}{
		{
			testName: "background context",
			ctx:      context.Background(),
		},
		{
			testName: "not user",
			ctx:      context.WithValue(context.Background(), UserContextKey, "zoop"),
		},
		{
			testName: "user is nil",
			ctx:      context.WithValue(context.Background(), UserContextKey, nilUser),
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("FetchUserFromContext - Errors - %v", given.testName), func(t *testing.T) {
			result, err := FetchUserFromContext(given.ctx)
			require.NotNil(t, err, "no error in FetchUserFromContext: %v", err)
			require.Empty(t, result, "Expected empty result, got: %v", result)
		})
	}
}

func TestFetchUserFromContextSuccesses(t *testing.T) {
	user := NewDummyUser(t)

	var tests = []struct {
		testName string
		ctx      context.Context
		expected *User
	}{
		{
			testName: "user is nil",
			ctx:      context.WithValue(context.Background(), UserContextKey, user),
			expected: user,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("FetchUserFromContext - Successes - %v", given.testName), func(t *testing.T) {
			result, err := FetchUserFromContext(given.ctx)
			require.Nil(t, err, "no error in FetchUserFromContext: %v", err)
			require.Equal(t, result, given.expected, "Result %v did not equal expected %v", result, given.expected)
		})
	}
}
