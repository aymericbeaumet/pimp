package sql

import (
	"database/sql"

	// https://github.com/golang/go/wiki/SQLDrivers
	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
	_ "modernc.org/sqlite"             // sqlite driver
)

// SQLOpen creates a new connection to an SQL database:
//   - SQLite: (SQLOpen "sqlite" "/tmp/pimp.db").Query "..."
func Open(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{
		driverName:     driverName,
		dataSourceName: dataSourceName,
		db:             db,
	}, nil
}
