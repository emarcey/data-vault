package database

import (
	"context"
	"fmt"

	"emarcey/data-vault/common"
)

func ListTablePermissions(ctx context.Context, db Database, userId string) ([]*common.TablePermission, error) {
	operation := "ListTablePermissions"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dtp.user_id,
			dtp.table_id,
			dt.name AS table_name,
			dtp.is_decrypt_allowed,
			dtp.created_by,
			dtp.updated_by
	FROM	admin.data_table_permissions dtp
	JOIN	admin.data_tables dt
		ON	dtp.table_id = dt.id
	WHERE	dtp.user_id = $1
		AND dtp.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query, userId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var tablePermissions []*common.TablePermission

	for rows.Next() {
		var row common.TablePermission
		err = rows.Scan(&row.UserId, &row.TableId, &row.TableName, &row.IsDecryptAllowed, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		tablePermissions = append(tablePermissions, &row)
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return tablePermissions, nil
}

func ListTablePermissionsForTable(ctx context.Context, db Database, tableId string) ([]*common.TablePermission, error) {
	operation := "ListTablePermissionsForTable"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dtp.user_id,
			dtp.table_id,
			dt.name AS table_name,
			dtp.is_decrypt_allowed,
			dtp.created_by,
			dtp.updated_by
	FROM	admin.data_table_permissions dtp
	JOIN	admin.data_tables dt
		ON	dtp.table_id = dt.id
	WHERE	dtp.table_id = $1
		AND dtp.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query, tableId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var tablePermissions []*common.TablePermission

	for rows.Next() {
		var row common.TablePermission
		err = rows.Scan(&row.UserId, &row.TableId, &row.TableName, &row.IsDecryptAllowed, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		tablePermissions = append(tablePermissions, &row)
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return tablePermissions, nil
}

func DeleteTablePermission(ctx context.Context, db Database, adminUser *common.User, userId, tableId string) error {
	operation := "DeleteTablePermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.data_table_permissions
	SET is_active = false, updated_by = $1
	WHERE	user_id = $2
		AND table_id = $3
	`
	result, err := db.ExecContext(tracer.Context(), query, adminUser.Id, userId, tableId)
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

func CreateTablePermission(ctx context.Context, db Database, adminUser *common.User, userId, tableId string, isDecryptAllowed bool) (*common.TablePermission, error) {
	operation := "CreateTablePermission"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.data_table_permissions (user_id, table_id, is_decrypt_allowed, is_active, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $5, $5)
	ON CONFLICT (user_id, table_id)
	DO UPDATE SET is_active = true, updated_by = $5
	RETURNING user_id, table_id, is_decrypt_allowed, is_active, created_by, updated_by
	`
	rows, err := db.QueryContext(tracer.Context(), query, userId, tableId, isDecryptAllowed, true, adminUser.Id)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}

	var tablePermission *common.TablePermission

	for rows.Next() {
		var row common.TablePermission
		err = rows.Scan(&row.UserId, &row.TableId, &row.TableName, &row.IsDecryptAllowed, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		tablePermission = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	if tablePermission == nil {
		return nil, common.NewResourceNotFoundError(operation, "table-permission", fmt.Sprintf("%v - %v", userId, tableId))
	}
	return tablePermission, nil
}
