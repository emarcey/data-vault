package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDeleteSecretPermissionErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(nil).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecretPermission - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecretPermission(context.Background(), dbMock, "callingUserId", "userId", "secretId")
			require.NotNil(t, err, "no error in DeleteSecretPermission: %v", err)
		})
	}
}

func TestDeleteSecretPermissionSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("callingUserId", "userId", "secretId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecretPermission - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecretPermission(context.Background(), dbMock, "callingUserId", "userId", "secretId")
			require.Nil(t, err, "error in DeleteSecretPermission: %v", err)
		})
	}
}

func TestCreateSecretPermissionErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(nil).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecretPermission - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecretPermission(context.Background(), dbMock, "callingUserId", "userId", "secretId")
			require.NotNil(t, err, "no error in CreateSecretPermission: %v", err)
		})
	}
}

func TestCreateSecretPermissionSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("userId", "secretId", "callingUserId", "callingUserId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecretPermission - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecretPermission(context.Background(), dbMock, "callingUserId", "userId", "secretId")
			require.Nil(t, err, "error in CreateSecretPermission: %v", err)
		})
	}
}
