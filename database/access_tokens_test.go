package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/emarcey/data-vault/common"
)

func TestSelectAccessTokensForAuthErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id_hash", "user_id", "invalid_at", "is_latest"}).AddRow("idHash", "userId", time.Now(), true).RowError(0, fmt.Errorf("oh no not the row"))).RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("SelectAccessTokensForAuth - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := SelectAccessTokensForAuth(context.Background(), dbMock)
			require.NotNil(t, err, "no error in SelectAccessTokensForAuth: %v", err)
			require.Empty(t, result, "Expected empty result, got: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestSelectAccessTokensForAuthSuccesses(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now()
	var inits = []struct {
		initFunc initFunc
		expected map[string]*common.AccessToken
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id_hash", "user_id", "invalid_at", "is_latest"})).RowsWillBeClosed()
			},
			expected: map[string]*common.AccessToken{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id_hash", "user_id", "invalid_at", "is_latest"}).AddRow("idHash1", "userId1", t1, true),
				).RowsWillBeClosed()
			},
			expected: map[string]*common.AccessToken{
				"idHash1": &common.AccessToken{
					Id:        "idHash1",
					UserId:    "userId1",
					InvalidAt: t1,
					IsLatest:  true,
				},
			},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows(
						[]string{"id_hash", "user_id", "invalid_at", "is_latest"}).
						AddRow("idHash1", "userId1", t1, true).
						AddRow("idHash2", "userId1", t2, true).
						AddRow("idHash2", "userId2", t2, true),
				).RowsWillBeClosed()
			},
			expected: map[string]*common.AccessToken{
				"idHash1": &common.AccessToken{
					Id:        "idHash1",
					UserId:    "userId1",
					InvalidAt: t1,
					IsLatest:  true,
				},
				"idHash2": &common.AccessToken{
					Id:        "idHash2",
					UserId:    "userId2",
					InvalidAt: t2,
					IsLatest:  true,
				},
			},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("SelectAccessTokensForAuth - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := SelectAccessTokensForAuth(context.Background(), dbMock)
			require.Nil(t, err, "Unexpected error in SelectAccessTokensForAuth: %v", err)
			require.Equal(t, result, given.expected, "Result, %+v, did not equal expected, %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeprecateLatestAccessTokenErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("UPDATE").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1").RowError(0, fmt.Errorf("oh no not the row"))).RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeprecateLatestAccessToken - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := DeprecateLatestAccessToken(context.Background(), dbMock, "userId")
			require.NotNil(t, err, "no error in DeprecateLatestAccessToken: %v", err)
			require.Empty(t, result, "Expected empty result, got: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeprecateLatestAccessTokenSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("UPDATE").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1")).RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeprecateLatestAccessToken - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := DeprecateLatestAccessToken(context.Background(), dbMock, "userId")
			require.Nil(t, err, "error in DeprecateLatestAccessToken: %v", err)
			require.NotEmpty(t, result, "Result was empty")
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateAccessTokenErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateAccessToken - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateAccessToken(context.Background(), dbMock, "userId", "accessTokenHash", time.Now())
			require.NotNil(t, err, "no error in CreateAccessToken: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateAccessTokenSuccesses(t *testing.T) {
	tmpTime := time.Now()
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("accessTokenHash", "userId", true, tmpTime)
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateAccessToken - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateAccessToken(context.Background(), dbMock, "userId", "accessTokenHash", tmpTime)
			require.Nil(t, err, "error in CreateAccessToken: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}
