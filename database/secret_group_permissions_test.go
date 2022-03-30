package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDeleteSecretGroupPermissionErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(nil).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecretGroupPermission - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecretGroupPermission(context.Background(), dbMock, "callingUserId", "userGroupId", "secretId")
			require.NotNil(t, err, "no error in DeleteSecretGroupPermission: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expecations not met: %v", err)
		})
	}
}

func TestDeleteSecretGroupPermissionSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("callingUserId", "userGroupId", "secretId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecretGroupPermission - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecretGroupPermission(context.Background(), dbMock, "callingUserId", "userGroupId", "secretId")
			require.Nil(t, err, "error in DeleteSecretGroupPermission: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expecations not met: %v", err)
		})
	}
}

func TestCreateSecretGroupPermissionErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(nil).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecretGroupPermission - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecretGroupPermission(context.Background(), dbMock, "callingUserId", "userGroupId", "secretId")
			require.NotNil(t, err, "no error in CreateSecretGroupPermission: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expecations not met: %v", err)
		})
	}
}

func TestCreateSecretGroupPermissionSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("userGroupId", "secretId", "callingUserId", "callingUserId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecretGroupPermission - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecretGroupPermission(context.Background(), dbMock, "callingUserId", "userGroupId", "secretId")
			require.Nil(t, err, "error in CreateSecretGroupPermission: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expecations not met: %v", err)
		})
	}
}
