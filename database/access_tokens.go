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
		tracer.CaptureException(err)
		return nil, common.NewDatabaseError(operation, err.Error())
	}
	defer rows.Close()

	accessTokenMap := make(map[string]*common.AccessToken)

	for rows.Next() {
		var row common.AccessToken
		err = rows.Scan(&row.Id, &row.UserId, &row.InvalidAt, &row.IsLatest)
		if err != nil {
			newErr := common.NewDatabaseError(operation, "Error in scan operation: %v", err)
			tracer.CaptureException(newErr)
			return nil, newErr
		}
		accessTokenMap[row.Id] = &row
	}
	err = rows.Err()
	if err != nil {
		newErr := common.NewDatabaseError(operation, "Error in rows.Err() operation: %v", err)
		tracer.CaptureException(newErr)
		return nil, newErr
	}
	return accessTokenMap, nil
}
