package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/tsmweb/user-service/config"
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
		config.DBHost(), config.DBPort(), config.DBUser(), config.DBPassword(), config.DBName(), config.DBSchema())
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return &PostgresDatabase{db}
}

// DB get instance of a connection to the database.
func (pd *PostgresDatabase) DB() *sql.DB {
	return pd.db
}
