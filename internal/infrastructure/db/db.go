package db

import (
	"database/sql"
	"fmt"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DB struct {
	*bun.DB
}

func NewConnection(conf *config.Config) (DB, error) {
	//dsn := "postgres://postgres:@localhost:5432/test?sslmode=disable"
	configDB := conf.Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", configDB.Username, configDB.Password, configDB.Host,
		configDB.Port, configDB.DbName)
	fmt.Printf("dsn: %s\n", dsn)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return DB{
		DB: db,
	}, nil
}
