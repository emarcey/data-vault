package database

import (
	"context"
	"time"

	"emarcey/data-vault/common"
)

func SelectAccessTokensForAuth(ctx context.Context, db *Database) (map[string]*common.AccessToken, error) {
	operation := "SelectAccessTokensForAuth"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	SELECT	at.id_hash,
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

func DeprecateLatestAccessToken(ctx context.Context, db *Database, userId string) error {
	operation := "DeprecateLatestAccessToken"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	UPDATE  admin.access_tokens
	SET 	is_latest = false
	WHERE	user_id = $1
		AND	is_latest = true
	`
	result, err := db.ExecContext(tracer.Context(), query, userId)
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
	db.logger.Debugf("%s updated %d rows", operation, rowsAffected)

	return nil
}

func CreateAccessToken(ctx context.Context, db *Database, userId, accessTokenHash string, invalidAt time.Time) error {
	operation := "CreateAccessToken"
	tracer := db.tracerCreator(ctx, operation)
	defer tracer.Close()

	query := `
	INSERT INTO  admin.access_tokens (id_hash, user_id, is_latest, invalid_at)
	VALUES($1, $2, $3, $4)
	`
	result, err := db.ExecContext(tracer.Context(), query, accessTokenHash, userId, true, invalidAt)
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
	db.logger.Debugf("%s created %d rows", operation, rowsAffected)

	return nil
}
