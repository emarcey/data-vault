package database

import (
	"context"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"emarcey/data-vault/common"
)

func TestCreateSecretErrors(t *testing.T) {
	secret1 := common.NewDummySecret(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecret - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecret(context.Background(), dbMock, secret1)
			require.NotNil(t, err, "no error in CreateSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestCreateSecretSuccesses(t *testing.T) {
	secret1 := common.NewDummySecret(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("CreateSecret - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = CreateSecret(context.Background(), dbMock, secret1)
			require.Nil(t, err, "error in CreateSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetSecretByNameErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	secret1 := common.NewDummySecret(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value", "description", "created_by", "updated_by"}).
				AddRow(secret1.Id, secret1.Name, secret1.Value, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy).
				RowError(0, fmt.Errorf("oh no not the row"))).RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetSecretByName - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := GetSecretByName(context.Background(), dbMock, user1, "secretName")
			require.NotNil(t, err, "no error in GetSecretByName: %v", err)
			require.Nil(t, result, "Expected nil result, got: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetSecretByNameSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	secret1 := common.NewDummySecret(t)

	var inits = []struct {
		initFunc initFunc
		expected *common.Secret
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value", "description", "created_by", "updated_by"}).
					AddRow(secret1.Id, secret1.Name, secret1.Value, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy)).RowsWillBeClosed()
			},
			expected: secret1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetSecretByName - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := GetSecretByName(context.Background(), dbMock, user1, "secretName")
			require.Nil(t, err, "Unexpected error in GetSecretByName: %v", err)
			require.Equal(t, result, given.expected, "Result %+v does not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetSecretIdWithWriteAccessErrors(t *testing.T) {
	user1 := common.NewDummyUser(t)
	secret1 := common.NewDummySecret(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value", "description", "created_by", "updated_by"}).
				AddRow(secret1.Id, secret1.Name, secret1.Value, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy).
				RowError(0, fmt.Errorf("oh no not the row"))).RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetSecretIdWithWriteAccess - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := GetSecretIdWithWriteAccess(context.Background(), dbMock, user1, "secretName")
			require.NotNil(t, err, "no error in GetSecretIdWithWriteAccess: %v", err)
			require.Empty(t, result, "Expected nil result, got: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestGetSecretIdWithWriteAccessSuccesses(t *testing.T) {
	user1 := common.NewDummyUser(t)
	secret1 := common.NewDummySecret(t)

	var inits = []struct {
		initFunc initFunc
		expected *common.Secret
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(secret1.Id)).RowsWillBeClosed()
			},
			expected: secret1,
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("GetSecretIdWithWriteAccess - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := GetSecretIdWithWriteAccess(context.Background(), dbMock, user1, "secretName")
			require.Nil(t, err, "Unexpected error in GetSecretIdWithWriteAccess: %v", err)
			require.Equal(t, result, given.expected.Id, "Result %+v does not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteSecretErrors(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("zoop")))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecret - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecret(context.Background(), dbMock, "userId", "secretName")
			require.NotNil(t, err, "no error in DeleteSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestDeleteSecretSuccesses(t *testing.T) {
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("DeleteSecret - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			err = DeleteSecret(context.Background(), dbMock, "userId", "secretName")
			require.Nil(t, err, "error in DeleteSecret: %v", err)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListSecretsErrors(t *testing.T) {
	secret1 := common.NewDummySecret(t)
	user1 := common.NewDummyUser(t)
	var inits = []initFunc{
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("Oh no!"))
		},
		func(dbMock *MockDatabase) {
			dbMock.mock.ExpectQuery("SELECT").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "updated_by"}).
					AddRow(secret1.Id, secret1.Name, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy).
					RowError(0, fmt.Errorf("oh no not the row"))).
				RowsWillBeClosed()
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListSecrets - Errors - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given(dbMock)

			result, err := ListSecrets(context.Background(), dbMock, user1, 0, 10)
			require.NotNil(t, err, "no error in ListSecrets: %v", err)
			require.Nil(t, result, "Result was not nil: %v", result)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}

func TestListSecretsSuccesses(t *testing.T) {
	secret1 := common.NewDummySecret(t)
	secret1.Value = ""
	secret2 := common.NewDummySecret(t)
	secret2.Value = ""
	user1 := common.NewDummyUser(t)
	var inits = []struct {
		initFunc initFunc
		expected []*common.Secret
	}{
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "updated_by"})).
					RowsWillBeClosed()
			},
			expected: []*common.Secret{},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "updated_by"}).
						AddRow(secret1.Id, secret1.Name, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy)).
					RowsWillBeClosed()
			},
			expected: []*common.Secret{secret1},
		},
		{
			initFunc: func(dbMock *MockDatabase) {
				dbMock.mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "updated_by"}).
						AddRow(secret1.Id, secret1.Name, secret1.Description, secret1.CreatedBy, secret1.UpdatedBy).
						AddRow(secret2.Id, secret2.Name, secret2.Description, secret2.CreatedBy, secret2.UpdatedBy)).
					RowsWillBeClosed()
			},
			expected: []*common.Secret{secret1, secret2},
		},
	}

	for idx, given := range inits {
		t.Run(fmt.Sprintf("ListSecrets - Successes - %v", idx), func(t *testing.T) {
			dbMock, err := NewMockDatabase()
			require.Nil(t, err, "Unexpected err creating mock db: %v", err)
			given.initFunc(dbMock)

			result, err := ListSecrets(context.Background(), dbMock, user1, 0, 10)
			require.Nil(t, err, "no error in ListSecrets: %v", err)
			require.Equal(t, result, given.expected, "Result %+v did not equal expected %+v", result, given.expected)
			err = dbMock.mock.ExpectationsWereMet()
			require.Nil(t, err, "expectations not met: %v", err)
		})
	}
}
