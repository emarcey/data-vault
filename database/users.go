package database

import (
	"context"

	"emarcey/data-vault/common"
)

func SelectUsersForAuth(ctx context.Context, db Database) (map[string]*common.User, error) {
	operation := "SelectUsersForAuth"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	u.id,
			u.name,
			u.is_active,
			u.type,
			u.client_secret_hash
	FROM	admin.users u
	WHERE	u.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	userMap := make(map[string]*common.User)

	for rows.Next() {
		var row common.User
		err = rows.Scan(&row.Id, &row.Name, &row.IsActive, &row.Type, &row.SecretHash)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		userMap[row.Id] = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return userMap, nil
}

func ListUsers(ctx context.Context, db Database) ([]*common.User, error) {
	operation := "ListUsers"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	u.id,
			u.name,
			u.is_active,
			u.type
	FROM	admin.users u
	WHERE	u.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var users []*common.User

	for rows.Next() {
		var row common.User
		err = rows.Scan(&row.Id, &row.Name, &row.IsActive, &row.Type)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		users = append(users, &row)
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return users, nil
}

func GetUserById(ctx context.Context, db Database, userId string) (*common.User, error) {
	operation := "GetUserById"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	u.id,
			u.name,
			u.is_active,
			u.type
	FROM	admin.users u
	WHERE	id = $1
		AND u.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query, userId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var user *common.User

	for rows.Next() {
		var row common.User
		err = rows.Scan(&row.Id, &row.Name, &row.IsActive, &row.Type)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		user = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	if user == nil {
		return nil, common.NewResourceNotFoundError(operation, "id", userId)
	}
	return user, nil
}

func DeleteUser(ctx context.Context, db Database, userId string) error {
	operation := "DeleteUser"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.users
	SET is_active = false
	WHERE	id = $1
	`
	result, err := db.ExecContext(tracer.Context(), query, userId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return dbErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return dbErr
	}
	db.GetLogger().Debugf("%s soft deleted %d rows", operation, rowsAffected)

	return nil
}

func CreateUser(ctx context.Context, db Database, userId, userName, userType, userSecretHash string) error {
	operation := "CreateUser"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.users (id, name, is_active, type, client_secret_hash)
	VALUES($1, $2, $3, $4, $5)
	`
	result, err := db.ExecContext(tracer.Context(), query, userId, userName, true, userType, userSecretHash)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return dbErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return dbErr
	}
	db.GetLogger().Debugf("%s created %d rows", operation, rowsAffected)

	return nil
}