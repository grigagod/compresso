// Package postgres provides connection to PosgresSQL.
package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// NewPsqlDB return new sqlx.DB instance.
func NewPsqlDB(host string, port string, user string, dbname string, password string, driver string) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbname,
		password,
	)

	db, err := sqlx.Connect(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
