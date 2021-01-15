package database

import (
	"database/sql"
	"fmt"
	"github.com/tsmweb/use-service/helper/setting"
	"time"
)

const (
	dbDriver = "postgres"
)

var (
	instance Database
)

// Database read only interface to access database connection.
type Database interface {
	DB() *sql.DB
}

// PostgresDatabase stores a reference to the bank connection pool.
type PostgresDatabase struct {
	db *sql.DB
}

// NewPostgresDatabase creates a new instance of PostgresDatabase.
func NewPostgresDatabase() Database {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s,public",
		setting.DBHost(), setting.DBPort(), setting.DBUser(), setting.DBPassword(), setting.DBName(), setting.DBSchema())
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute * 5)

	return &PostgresDatabase{db}
}

// DB get instance of a connection to the database.
func (pd *PostgresDatabase) DB() *sql.DB {
	return pd.db
}
