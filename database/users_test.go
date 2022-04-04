package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"emarcey/data-vault/common"
)

func TestSelectUsersForAuthErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type", "client_secret_hash"}).
					AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type, user1.SecretHash).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("SelectUsersForAuth - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := SelectUsersForAuth(context.Background(), dbMock)
			require.NotNil(t, err, "no error in SelectUsersForAuth: %v", err)
			require.Empty(t, result, "Expected empty result, got: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestSelectUsersForAuthSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	user2 := common.NewDummyUser(t)
	user3 := common.NewDummyUser(t)
	var inits = []struct {
		initFunc initFunc
		expected map[string]*common.User
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type", "client_secret_hash"})).
					RowsWillBeClosed()
			},
			expected: map[string]*common.User{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "is_active", "type", "client_secret_hash"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type, user1.SecretHash),
				).RowsWillBeClosed()
			},
			expected: map[string]*common.User{
				user1.Id: user1,
			},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "name", "is_active", "type", "client_secret_hash"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type, user1.SecretHash).
						AddRow(user2.Id, user2.Name, user2.IsActive, user2.Type, user2.SecretHash).
						AddRow(user3.Id, user3.Name, user3.IsActive, user3.Type, user3.SecretHash),
				).RowsWillBeClosed()
			},
			expected: map[string]*common.User{
				user1.Id: user1,
				user2.Id: user2,
				user3.Id: user3,
			},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("SelectUsersForAuth - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := SelectUsersForAuth(context.Background(), dbMock)
			require.Nil(t, err, "Unexpected error in SelectUsersForAuth: %v", err)
			require.Equal(t, result, given.expected, "Result, %+v, did not equal expected, %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetUserByIdErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(user1.Id, user1.Name).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetUserById - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := GetUserById(context.Background(), dbMock, "userId")
			require.NotNil(t, err, "no error in GetUserById: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetUserByIdSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []struct {
		initFunc initFunc
		expected *common.User
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WithArgs("userId").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type)).
					RowsWillBeClosed()
			},
			expected: user1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetUserById - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := GetUserById(context.Background(), dbMock, "userId")
			require.Nil(t, err, "no error in GetUserById: %v", err)
			given.expected.SecretHash = ""
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteUserErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUser - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUser(context.Background(), dbMock, "callingUserId", "userId")
			require.NotNil(t, err, "no error in DeleteUser: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteUserSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("callingUserId", "userId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUser - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUser(context.Background(), dbMock, "callingUserId", "userId")
			require.Nil(t, err, "error in DeleteUser: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("INSERT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(user1.Id, user1.Name).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUser - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := CreateUser(context.Background(), dbMock, "callingUserId", user1.Id, user1.Name, user1.Type, user1.SecretHash)
			require.NotNil(t, err, "no error in CreateUser: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []struct {
		initFunc initFunc
		expected *common.User
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("INSERT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type)).
					RowsWillBeClosed()
			},
			expected: user1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUser - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := CreateUser(context.Background(), dbMock, "callingUserId", user1.Id, user1.Name, user1.Type, user1.SecretHash)
			require.Nil(t, err, "no error in CreateUser: %v", err)
			given.expected.SecretHash = ""
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestRotateUserSecretErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("RotateUserSecret - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = RotateUserSecret(context.Background(), dbMock, "userId", "secretHash")
			require.NotNil(t, err, "no error in RotateUserSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestRotateUserSecretSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("secretHash", "userId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("RotateUserSecret - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = RotateUserSecret(context.Background(), dbMock, "userId", "secretHash")
			require.Nil(t, err, "error in RotateUserSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUsersErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
					AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUsers - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := ListUsers(context.Background(), dbMock, 0, 10)
			require.NotNil(t, err, "no error in ListUsers: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUsersSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	user1.SecretHash = ""
	user2 := common.NewDummyUser(t)
	user2.SecretHash = ""
	var inits = []struct {
		initFunc initFunc
		expected []*common.User
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"})).
					RowsWillBeClosed()
			},
			expected: []*common.User{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type)).
					RowsWillBeClosed()
			},
			expected: []*common.User{user1},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type).
						AddRow(user2.Id, user2.Name, user2.IsActive, user2.Type)).
					RowsWillBeClosed()
			},
			expected: []*common.User{user1, user2},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUsers - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := ListUsers(context.Background(), dbMock, 0, 10)
			require.Nil(t, err, "no error in ListUsers: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUsersInGroupErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
					AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUsersInGroup - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := ListUsersInGroup(context.Background(), dbMock, "userGroupId1", 0, 10)
			require.NotNil(t, err, "no error in ListUsersInGroup: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUsersInGroupSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	user1.SecretHash = ""
	user2 := common.NewDummyUser(t)
	user2.SecretHash = ""
	var inits = []struct {
		initFunc initFunc
		expected []*common.User
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"})).
					RowsWillBeClosed()
			},
			expected: []*common.User{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type)).
					RowsWillBeClosed()
			},
			expected: []*common.User{user1},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active", "type"}).
						AddRow(user1.Id, user1.Name, user1.IsActive, user1.Type).
						AddRow(user2.Id, user2.Name, user2.IsActive, user2.Type)).
					RowsWillBeClosed()
			},
			expected: []*common.User{user1, user2},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUsersInGroup - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := ListUsersInGroup(context.Background(), dbMock, "userGroupId1", 0, 10)
			require.Nil(t, err, "no error in ListUsersInGroup: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}
