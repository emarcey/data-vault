package database

import (
	"context"

	"github.com/emarcey/data-vault/common"
)

func DeleteUserGroupMember(ctx context.Context, db Database, callingUserId, userGroupId, userId string) error {
	operation := "DeleteUserGroupMember"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.user_group_members
	SET is_active = false,
		updated_by = $1
	WHERE	user_group_id = $2 and user_id = $3
	`
	result, err := db.ExecContext(tracer.Context(), query, callingUserId, userGroupId, userId)
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

func CreateUserGroupMember(ctx context.Context, db Database, callingUserId, userGroupId, userId string) error {
	operation := "CreateUserGroupMember"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.user_group_members (user_group_id, user_id, created_by, updated_by)
	VALUES($1, $2, $3, $4)
	`
	result, err := db.ExecContext(tracer.Context(), query, userGroupId, userId, callingUserId, callingUserId)
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
