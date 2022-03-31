package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDeleteUserGroupMemberErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUserGroupMember - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUserGroupMember(context.Background(), dbMock, "callingUserId", "userGroupId", "userId")
			require.NotNil(t, err, "no error in DeleteUserGroupMember: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteUserGroupMemberSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("callingUserId", "userGroupId", "userId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUserGroupMember - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUserGroupMember(context.Background(), dbMock, "callingUserId", "userGroupId", "userId")
			require.Nil(t, err, "error in DeleteUserGroupMember: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserGroupMemberErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUserGroupMember - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateUserGroupMember(context.Background(), dbMock, "callingUserId", "userGroupId", "userId")
			require.NotNil(t, err, "no error in CreateUserGroupMember: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserGroupMemberSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("userGroupId", "userId", "callingUserId", "callingUserId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUserGroupMember - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateUserGroupMember(context.Background(), dbMock, "callingUserId", "userGroupId", "userId")
			require.Nil(t, err, "error in CreateUserGroupMember: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}
