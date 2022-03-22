package database

import (
	"context"

	"emarcey/data-vault/common"
)

func CreateSecret(ctx context.Context, db Database, secret *common.Secret) error {
	operation := "CreateSecret"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.secrets (id, name, value, description, created_by, updated_by)
	VALUES($1, $2, $3, $4, $5, $6)
	`
	result, err := db.ExecContext(tracer.Context(), query, secret.Id, secret.Name, secret.Value, secret.Description, secret.CreatedBy, secret.UpdatedBy)
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
