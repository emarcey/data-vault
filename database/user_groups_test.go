package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/emarcey/data-vault/common"
)

func TestDeleteUserGroupErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUserGroup - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUserGroup(context.Background(), dbMock, "callingUserId", "userGroupId")
			require.NotNil(t, err, "no error in DeleteUserGroup: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteUserGroupSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs("callingUserId", "userGroupId")
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteUserGroup - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteUserGroup(context.Background(), dbMock, "callingUserId", "userGroupId")
			require.Nil(t, err, "error in DeleteUserGroup: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserGroupErrors(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("INSERT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(userGroup1.Id, userGroup1.Name).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUserGroup - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := CreateUserGroup(context.Background(), dbMock, "callingUserId", userGroup1.Id, userGroup1.Name)
			require.NotNil(t, err, "no error in CreateUserGroup: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateUserGroupSuccesses(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	var inits = []struct {
		initFunc initFunc
		expected *common.UserGroup
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("INSERT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(userGroup1.Id, userGroup1.Name)).
					RowsWillBeClosed()
			},
			expected: userGroup1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateUserGroup - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := CreateUserGroup(context.Background(), dbMock, "callingUserId", userGroup1.Id, userGroup1.Name)
			require.Nil(t, err, "no error in CreateUserGroup: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetUserGroupErrors(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(userGroup1.Id, userGroup1.Name).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetUserGroup - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := GetUserGroup(context.Background(), dbMock, userGroup1.Id)
			require.NotNil(t, err, "no error in GetUserGroup: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetUserGroupSuccesses(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	var inits = []struct {
		initFunc initFunc
		expected *common.UserGroup
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(userGroup1.Id, userGroup1.Name)).
					RowsWillBeClosed()
			},
			expected: userGroup1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetUserGroup - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := GetUserGroup(context.Background(), dbMock, userGroup1.Id)
			require.Nil(t, err, "no error in GetUserGroup: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUserGroupsErrors(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(userGroup1.Id, userGroup1.Name).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUserGroups - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := ListUserGroups(context.Background(), dbMock, 0, 10)
			require.NotNil(t, err, "no error in ListUserGroups: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListUserGroupsSuccesses(t *testing.T) {
	userGroup1 := common.NewDummyUserGroup(t)
	userGroup2 := common.NewDummyUserGroup(t)
	var inits = []struct {
		initFunc initFunc
		expected []*common.UserGroup
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"})).
					RowsWillBeClosed()
			},
			expected: []*common.UserGroup{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(userGroup1.Id, userGroup1.Name)).
					RowsWillBeClosed()
			},
			expected: []*common.UserGroup{userGroup1},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(userGroup1.Id, userGroup1.Name).
						AddRow(userGroup2.Id, userGroup2.Name)).
					RowsWillBeClosed()
			},
			expected: []*common.UserGroup{userGroup1, userGroup2},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListUserGroups - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := ListUserGroups(context.Background(), dbMock, 0, 10)
			require.Nil(t, err, "no error in ListUserGroups: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}
