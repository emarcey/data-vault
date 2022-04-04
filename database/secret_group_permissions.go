package database

import (
	"context"

	"github.com/emarcey/data-vault/common"
)

func DeleteSecretGroupPermission(ctx context.Context, db Database, callingUserId, userGroupId, secretId string) error {
	operation := "DeleteSecretGroupPermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.secret_group_permissions
	SET is_active = false,
		updated_by = $1
	WHERE	user_group_id = $2 and secret_id = $3
	`

	result, err := db.ExecContext(tracer.Context(), query, callingUserId, userGroupId, secretId)
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

func CreateSecretGroupPermission(ctx context.Context, db Database, callingUserId, userGroupId, secretId string) error {
	operation := "CreateSecretGroupPermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.secret_group_permissions (user_group_id, secret_id, created_by, updated_by)
	VALUES($1, $2, $3, $4)
	`
	result, err := db.ExecContext(tracer.Context(), query, userGroupId, secretId, callingUserId, callingUserId)
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
