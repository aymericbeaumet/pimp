package sql

import (
	"context"
	"database/sql"
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"SQLOpen": Open,
	}
}

type Database struct {
	driverName     string
	dataSourceName string
	db             *sql.DB
}

func (db Database) String() string {
	return db.driverName
}

func (db Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(context.Background(), query, args...)
}

func (db Database) Ping() (interface{}, error) {
	return nil, db.db.PingContext(context.Background())
}

func (db Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(context.Background(), query, args...)
}
