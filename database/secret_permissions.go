package database

import (
	"context"

	"github.com/emarcey/data-vault/common"
)

func DeleteSecretPermission(ctx context.Context, db Database, callingUserId, userId, secretId string) error {
	operation := "DeleteSecretPermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.secret_permissions
	SET is_active = false,
		updated_by = $1
	WHERE	user_id = $2 and secret_id = $3
	`
	result, err := db.ExecContext(tracer.Context(), query, callingUserId, userId, secretId)
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

func CreateSecretPermission(ctx context.Context, db Database, callingUserId, userId, secretId string) error {
	operation := "CreateSecretPermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.secret_permissions (user_id, secret_id, created_by, updated_by)
	VALUES($1, $2, $3, $4)
	`
	result, err := db.ExecContext(tracer.Context(), query, userId, secretId, callingUserId, callingUserId)
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
