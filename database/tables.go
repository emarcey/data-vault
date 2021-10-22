package database

import (
	"context"

	"emarcey/data-vault/common"
)

func ListTables(ctx context.Context, db Database, user *common.User) ([]*common.Table, error) {
	operation := "ListTables"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dt.id,
			dt.name,
			dt.description,
			dt.created_by,
			dt.updated_by
	FROM	admin.data_tables dt
	LEFT JOIN	admin.data_table_permissions dtp
		ON	dtp.table_id = dt.id
		AND dtp.is_active
		AND dtp.user_id = $1
	WHERE	dt.is_active
		AND (dtp.table_id IS NOT NULL OR $2 = 'admin' OR dt.created_by = $3)
	`
	rows, err := db.QueryContext(tracer.Context(), query, user.Id, user.Type, user.Id)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var tables []*common.Table

	for rows.Next() {
		var row common.Table
		err = rows.Scan(&row.Id, &row.Name, &row.Description, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		tables = append(tables, &row)
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return tables, nil
}

func GetTableById(ctx context.Context, db Database, user *common.User, tableId string) (*common.Table, error) {
	operation := "GetTableById"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dt.id,
			dt.name,
			dt.description,
			dt.created_by,
			dt.updated_by
	FROM	admin.data_tables dt
	LEFT JOIN	admin.data_table_permissions dtp
		ON	dtp.table_id = dt.id
		AND dtp.is_active
		AND dtp.user_id = $1
	WHERE	dt.is_active
		AND dt.id = $2
		AND (dtp.table_id IS NOT NULL OR $3 = 'admin' OR dt.created_by = $4)
	`
	rows, err := db.QueryContext(tracer.Context(), query, user.Id, tableId, user.Type, user.Id)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var table *common.Table

	for rows.Next() {
		var row common.Table
		err = rows.Scan(&row.Id, &row.Name, &row.Description, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		table = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	if table == nil {
		return nil, common.NewResourceNotFoundError(operation, "id", tableId)
	}
	return table, nil
}

func DeleteTable(ctx context.Context, db Database, user *common.User, tableId string) error {
	operation := "DeleteTable"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
		UPDATE	admin.data_tables dt
		SET		is_active = false,
				updated_by = $1
		WHERE dt.id = $2
		`
	result, err := db.ExecContext(tracer.Context(), query, user.Id, tableId)
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

func CreateTable(ctx context.Context, db Database, user *common.User, tableName, tableDescription string) (*common.Table, error) {
	operation := "CreateTable"
	tracer := db.CreateTrace(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.data_tables (name, description, is_active, created_by, updated_by)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, name, description, created_by, updated_by
	`
	rows, err := db.QueryContext(tracer.Context(), query, tableName, tableDescription, true, user.Id, user.Id)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}

	var table *common.Table

	for rows.Next() {
		var row common.Table
		err = rows.Scan(&row.Id, &row.Name, &row.Description, &row.CreatedBy, &row.UpdatedBy)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		table = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	if table == nil {
		return nil, common.NewResourceNotFoundError(operation, "name", tableName)
	}
	return table, nil
}
