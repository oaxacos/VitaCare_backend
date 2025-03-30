package utils

import (
	"context"
	"fmt"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/uptrace/bun"
)

func CleanUpDB(db *bun.DB, ctx context.Context, tables []string) error {
	// clean up database
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.GetGlobalLogger().Error(err)
		return err
	}
	for _, table := range tables {
		_, err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			logger.GetGlobalLogger().Error(err)
			return tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.GetGlobalLogger().Error(err)
		return err
	}

	return nil
}
