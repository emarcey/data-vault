package database

import (
	"context"
	"fmt"

	"emarcey/data-vault/common"
)

func SelectUsersForAuth(ctx context.Context, db *Database) (map[string]*common.User, error) {
	operation := "SelectUsersForAuth"
	tracer := db.tracerCreator(ctx, operation)
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
		tracer.CaptureException(err)
		return nil, common.NewDatabaseError(operation, err.Error())
	}
	defer rows.Close()

	userMap := make(map[string]*common.User)

	for rows.Next() {
		var row common.User
		var clientSecretHash string
		err = rows.Scan(&row.Id, &row.Name, &row.IsActive, &row.Type, &clientSecretHash)
		if err != nil {
			newErr := common.NewDatabaseError(operation, "Error in scan operation: %v", err)
			tracer.CaptureException(newErr)
			return nil, newErr
		}
		userMap[fmt.Sprintf("%s_%s", row.Id, clientSecretHash)] = &row
	}
	err = rows.Err()
	if err != nil {
		newErr := common.NewDatabaseError(operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(newErr)
		return nil, newErr
	}
	return userMap, nil
}

func ListUsers(ctx context.Context, db *Database) ([]*common.User, error) {
	operation := "ListUsers"
	tracer := db.tracerCreator(ctx, operation)
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
		tracer.CaptureException(err)
		return nil, common.NewDatabaseError(operation, err.Error())
	}
	defer rows.Close()

	var users []*common.User

	for rows.Next() {
		var row common.User
		err = rows.Scan(&row.Id, &row.Name, &row.IsActive, &row.Type)
		if err != nil {
			newErr := common.NewDatabaseError(operation, "Error in scan operation: %v", err)
			tracer.CaptureException(newErr)
			return nil, newErr
		}
		users = append(users, &row)
	}
	err = rows.Err()
	if err != nil {
		newErr := common.NewDatabaseError(operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(newErr)
		return nil, newErr
	}
	return users, nil
}
