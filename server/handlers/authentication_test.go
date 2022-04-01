package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"emarcey/data-vault/common"
	"emarcey/data-vault/common/tracer"
	"emarcey/data-vault/dependencies"
)

var devUser = &common.User{
	Id:         "devUser",
	Name:       "devUser",
	Type:       "developer",
	IsActive:   true,
	SecretHash: "sha256:56112d265d7e45f567b3c8fd38c72844cba1224988bea814795a42b30278b623",
}
var adminUser = &common.User{
	Id:         "adminUser",
	Name:       "adminUser",
	Type:       "admin",
	IsActive:   true,
	SecretHash: "sha256:033a37715490c72ac56948f49595073a29f6aec382493ce7da48d04462bf5c70",
}

var userMap = map[string]*common.User{
	devUser.Id:   devUser,
	adminUser.Id: adminUser,
}

var testLogger = logrus.New()
var userCache = dependencies.NewMockUserCache(testLogger, userMap)

var devAccessToken = &common.AccessToken{
	Id:        "devAccessToken",
	UserId:    devUser.Id,
	InvalidAt: time.Now(),
	IsLatest:  true,
}
var adminAccessToken = &common.AccessToken{
	Id:        "adminAccessToken",
	UserId:    adminUser.Id,
	InvalidAt: time.Now(),
	IsLatest:  true,
}
var hangingAccessToken = &common.AccessToken{
	Id:        "hangingAccessToken",
	UserId:    "noUser",
	InvalidAt: time.Now(),
	IsLatest:  true,
}

var accessTokenMap = map[string]*common.AccessToken{
	"sha256:3f62e60c220787b3bb37b0d1d4987531135d44e8b5ad711569782c73487cf530": devAccessToken,
	"sha256:37a57e8270331014fa8b112862687344fe046a739af7872d615c2a209476d262": adminAccessToken,
	"sha256:10b0a44425a8b97a77a03084ee62c888a0f0253f2019df3ab78972dc37a3614c": hangingAccessToken,
}
var accessTokenCache = dependencies.NewMockAccessTokenCache(testLogger, accessTokenMap)

func TestAuthenticateClientErrors(t *testing.T) {
	var tests = []struct {
		testName   string
		ctx        context.Context
		checkAdmin bool
	}{
		{
			testName: "no client id",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"dummy": []string{},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "no client secret",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id": []string{"hi there"},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "user not found",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{"hi there"},
					"Client-Secret": []string{"user"},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "invalid secret",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{devUser.Id},
					"Client-Secret": []string{"user"},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "user not admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{devUser.Id},
					"Client-Secret": []string{devUser.Id},
				},
			}),
			checkAdmin: true,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("authenticateClient - Errors - %v", given.testName), func(t *testing.T) {
			result, err := authenticateClient(given.ctx, "op", tracer.NewNoOpTracer(given.ctx), userCache, given.checkAdmin)
			require.NotNil(t, err, "no error in authenticateClient: %v", err)
			require.Nil(t, result, "Expected empty result, got: %v", result)
		})
	}
}

func TestAuthenticateClientSuccesses(t *testing.T) {
	var tests = []struct {
		testName   string
		ctx        context.Context
		checkAdmin bool
		expected   *common.User
	}{
		{
			testName: "dev user - not checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{devUser.Id},
					"Client-Secret": []string{devUser.Id},
				},
			}),
			checkAdmin: false,
			expected:   devUser,
		},
		{
			testName: "admin user - not checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{adminUser.Id},
					"Client-Secret": []string{adminUser.Id},
				},
			}),
			checkAdmin: false,
			expected:   adminUser,
		},
		{
			testName: "admin user - checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Client-Id":     []string{adminUser.Id},
					"Client-Secret": []string{adminUser.Id},
				},
			}),
			checkAdmin: true,
			expected:   adminUser,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("authenticateClient - Successes - %v", given.testName), func(t *testing.T) {
			result, err := authenticateClient(given.ctx, "op", tracer.NewNoOpTracer(given.ctx), userCache, given.checkAdmin)
			require.Nil(t, err, "no error in authenticateClient: %v", err)
			require.Equal(t, result, given.expected, "Result %v did not equal expected %v", result, given.expected)
		})
	}
}

func TestAuthenticateAccessTokenErrors(t *testing.T) {
	var tests = []struct {
		testName   string
		ctx        context.Context
		checkAdmin bool
	}{
		{
			testName: "no access token",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"dummy": []string{},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "access token not found",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"hi there"},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "user not found",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"hangingAccessToken"},
				},
			}),
			checkAdmin: false,
		},
		{
			testName: "user not admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"devAccessToken"},
				},
			}),
			checkAdmin: true,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("authenticateAccessToken - Errors - %v", given.testName), func(t *testing.T) {
			result, err := authenticateAccessToken(given.ctx, "op", tracer.NewNoOpTracer(given.ctx), userCache, accessTokenCache, given.checkAdmin)
			require.NotNil(t, err, "no error in authenticateAccessToken: %v", err)
			require.Nil(t, result, "Expected empty result, got: %v", result)
		})
	}
}

func TestAuthenticateAccessTokenSuccesses(t *testing.T) {
	var tests = []struct {
		testName   string
		ctx        context.Context
		checkAdmin bool
		expected   *common.User
	}{
		{
			testName: "dev user - not checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"devAccessToken"},
				},
			}),
			checkAdmin: false,
			expected:   devUser,
		},
		{
			testName: "admin user - not checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"adminAccessToken"},
				},
			}),
			checkAdmin: false,
			expected:   adminUser,
		},
		{
			testName: "admin user - checking admin",
			ctx: common.InjectHeaderIntoContext(context.Background(), &http.Request{
				Header: map[string][]string{
					"Access-Token": []string{"adminAccessToken"},
				},
			}),
			checkAdmin: true,
			expected:   adminUser,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("authenticateAccessToken - Successes - %v", given.testName), func(t *testing.T) {
			result, err := authenticateAccessToken(given.ctx, "op", tracer.NewNoOpTracer(given.ctx), userCache, accessTokenCache, given.checkAdmin)
			require.Nil(t, err, "no error in authenticateAccessToken: %v", err)
			require.Equal(t, result, given.expected, "Result %v did not equal expected %v", result, given.expected)
		})
	}
}
