package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oaxacos/vitacare/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBRepository struct {
	*bun.DB
}

func NewConnection(conf *config.Config) (*DBRepository, error) {
	configDB := conf.Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", configDB.Username, configDB.Password, configDB.Host,
		configDB.Port, configDB.DbName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DBRepository{
		DB: db,
	}, nil
}

func (db *DBRepository) Close() error {
	return db.DB.Close()
}

func (db *DBRepository) WithTransaction(context context.Context, fn func(*bun.Tx) error) error {
	tx, err := db.DB.BeginTx(context, nil)
	if err != nil {
		return err
	}

	if err := fn(&tx); err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
