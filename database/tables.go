package database

import (
	"context"

	"emarcey/data-vault/common"
)

func ListTables(ctx context.Context, db *Database, user *common.User) ([]*common.Table, error) {
	operation := "ListTables"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dt.id,
			dt.name,
			dt.created_by,
			dt.updated_by
	FROM	admin.data_tables dt
	JOIN	admin.data_table_permissions dtp
		ON	dtp.table_id = dt.id
		AND dtp.is_active
		AND dtp.user_id = $1 or $2 = 'admin'
	WHERE	dt.is_active
	`
	rows, err := db.QueryContext(tracer.Context(), query, user.Id, user.Type)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var tables []*common.Table

	for rows.Next() {
		var row common.Table
		err = rows.Scan(&row.Id, &row.Name, &row.CreatedBy, &row.UpdatedBy)
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

func GetTableById(ctx context.Context, db *Database, user *common.User, tableId string) (*common.Table, error) {
	operation := "GetTableById"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	dt.id,
			dt.name,
			dt.created_by,
			dt.updated_by
	FROM	admin.data_tables dt
	JOIN	admin.data_table_permissions dtp
		ON	dtp.table_id = dt.id
		AND dtp.is_active
		AND dtp.user_id = $1 or $2 = 'admin'
	WHERE	dt.is_active
		AND dtp.id = $3
	`
	rows, err := db.QueryContext(tracer.Context(), query, user.Id, user.Type, tableId)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	var table *common.Table

	for rows.Next() {
		var row common.Table
		err = rows.Scan(&row.Id, &row.Name, &row.CreatedBy, &row.UpdatedBy)
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

func DeleteTable(ctx context.Context, db *Database, user *common.User, tableId string) error {
	operation := "DeleteTable"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE	admin.data_tables dt
	SET		is_active = False,
			updated_by = $1
	FROM	admin.data_table_permissions dtp
	WHERE	dtp.table_id = dt.id
		AND dtp.is_active
		AND dtp.user_id = $2 or $3 = 'admin'
		AND	dt.is_active
		AND dtp.id = $4
	`
	result, err := db.ExecContext(tracer.Context(), query, user.Id, user.Id, user.Type, tableId)
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
	db.logger.Debugf("%s soft deleted %d rows", operation, rowsAffected)
	return nil
}

// func DeleteTable(ctx context.Context, db *Database, tableName, userId string) error {
// 	operation := "DeleteTable"
// 	tracer := db.tracerCreator(ctx, operation)
// 	defer tracer.Close()

// 	query := `
// 	UPDATE  admin.data_tables
// 	SET is_active = false,
// 		updated_by = $1
// 	WHERE	id = $2
// 	`
// 	result, err := db.ExecContext(tracer.Context(), query, tableName, userId)
// 	if err != nil {
// 		dbErr := common.NewDatabaseError(err, operation, "")
// 		tracer.CaptureException(dbErr)
// 		return dbErr
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		dbErr := common.NewDatabaseError(err, operation, "")
// 		tracer.CaptureException(dbErr)
// 		return dbErr
// 	}
// 	db.logger.Debugf("%s updated %d rows", operation, rowsAffected)

// 	return nil
// }

// func CreateTable(ctx context.Context, db *Database, userId, userName, userType, userSecretHash string) error {
// 	operation := "CreateTable"
// 	tracer := db.tracerCreator(ctx, operation)
// 	defer tracer.Close()

// 	query := `
// 	INSERT INTO  admin.users (id, name, is_active, type, client_secret_hash)
// 	VALUES($1, $2, $3, $4, $5)
// 	`
// 	result, err := db.ExecContext(tracer.Context(), query, userId, userName, true, userType, userSecretHash)
// 	if err != nil {
// 		dbErr := common.NewDatabaseError(err, operation, "")
// 		tracer.CaptureException(dbErr)
// 		return dbErr
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		dbErr := common.NewDatabaseError(err, operation, "")
// 		tracer.CaptureException(dbErr)
// 		return dbErr
// 	}
// 	db.logger.Debugf("%s created %d rows", operation, rowsAffected)

// 	return nil
// }
