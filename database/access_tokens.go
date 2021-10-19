package database

import (
	"context"

	"emarcey/data-vault/common"
)

func SelectAccessTokensForAuth(ctx context.Context, db *Database) (map[string]*common.AccessToken, error) {
	operation := "SelectAccessTokensForAuth"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	at.id,
			at.user_id,
			at.invalid_at,
			at.is_latest
	FROM	admin.access_tokens at
	WHERE	at.is_latest
		AND at.invalid_at > NOW()
	`
	rows, err := db.QueryContext(tracer.Context(), query)
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "")
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	defer rows.Close()

	accessTokenMap := make(map[string]*common.AccessToken)

	for rows.Next() {
		var row common.AccessToken
		err = rows.Scan(&row.Id, &row.UserId, &row.InvalidAt, &row.IsLatest)
		if err != nil {
			dbErr := common.NewDatabaseError(err, operation, "Error in scan operation: %v", err)
			tracer.CaptureException(dbErr)
			return nil, dbErr
		}
		accessTokenMap[row.Id] = &row
	}
	err = rows.Err()
	if err != nil {
		dbErr := common.NewDatabaseError(err, operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(dbErr)
		return nil, dbErr
	}
	return accessTokenMap, nil
}
