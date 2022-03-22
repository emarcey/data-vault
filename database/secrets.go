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

func GetSecretByName(ctx context.Context, db Database, userId string) (*common.Secret, error) {
	operation := "GetSecretByName"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	s.id,
			s.name,
			s.value,
			s.description,
			created_by_user.name AS created_by,
			updated_by_user.name AS updated_by
	FROM	admin.secrets s
	JOIN	admin.users created_by_user
		ON 	s.created_by = created_by_user.id
		JOIN	admin.users updated_by_user
		ON 	s.updated_by = updated_by_user.id
	WHERE	s.name = $1
		AND s.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query, userId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var user *common.Secret

	for rows.Next() {
		var row common.Secret
		err = rows.Scan(&row.Id, &row.Name, &row.Value, &row.Description, &row.CreatedBy, &row.UpdatedBy)
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
