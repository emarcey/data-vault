package database

import (
	"context"

	"emarcey/data-vault/common"
)

func CreateUserGroup(ctx context.Context, db Database, callingUserId, userGroupId, userGroupName string) (*common.UserGroup, error) {
	operation := "CreateUserGroup"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.user_groups (id, name, is_active, created_by, updated_by)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, name
	`
	var userGroup *common.UserGroup
	rows, err := db.QueryContext(tracer.Context(), query, userGroupId, userGroupName, true, callingUserId, callingUserId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()
	for rows.Next() {
		var row common.UserGroup
		err = rows.Scan(&row.Id, &row.Name)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		userGroup = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	if userGroup == nil {
		return nil, common.NewResourceNotFoundError(operation, "id", userGroupId)
	}

	db.GetLogger().Debugf("%s created 1 row", operation)
	return userGroup, nil
}

func ListUserGroups(ctx context.Context, db Database) ([]*common.UserGroup, error) {
	operation := "ListUserGroups"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	u.id,
			u.name
	FROM	admin.user_groups u
	WHERE	u.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	userGroups := make([]*common.UserGroup, 0)

	for rows.Next() {
		var row common.UserGroup
		err = rows.Scan(&row.Id, &row.Name)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		userGroups = append(userGroups, &row)
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return userGroups, nil
}

func GetUserGroup(ctx context.Context, db Database, userId string) (*common.UserGroup, error) {
	operation := "GetUserGroup"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	u.id,
			u.name
	FROM	admin.user_groups u
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

	var user *common.UserGroup

	for rows.Next() {
		var row common.UserGroup
		err = rows.Scan(&row.Id, &row.Name)
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

func DeleteUserGroup(ctx context.Context, db Database, callingUserId, userGroupId string) error {
	operation := "DeleteUserGroup"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.user_groups
	SET is_active = false
		updated_by = $1
	WHERE	id = $2
	`
	result, err := db.ExecContext(tracer.Context(), query, callingUserId, userGroupId)
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
