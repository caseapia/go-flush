package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type Databases struct {
	Main *bun.DB
	Logs *bun.DB
}

func Connect(dbName string, maxOpen, maxIdle int) (*bun.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		dbName,
	)

	sqlDB, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return bun.NewDB(sqlDB, mysqldialect.New()), nil
}

func NewDatabases() (*Databases, error) {
	mainDB, err := Connect(os.Getenv("DB_MAIN_NAME"), 25, 10)
	if err != nil {
		return nil, err
	}

	logsDB, err := Connect(os.Getenv("DB_LOGS_NAME"), 5, 2)
	if err != nil {
		return nil, err
	}

	return &Databases{
		Main: mainDB,
		Logs: logsDB,
	}, nil
}
