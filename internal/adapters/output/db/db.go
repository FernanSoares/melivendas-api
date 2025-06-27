package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fesbarbosa/melivendas-api/internal/config"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {

	dsn := cfg.GetDSN()
	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {

		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1049 {

			dbConfigWithoutDB := *cfg
			dbConfigWithoutDB.DBName = ""
			rootDSN := dbConfigWithoutDB.GetDSN()

			rootDB, err := sqlx.Open(cfg.Driver, rootDSN)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
			}
			defer rootDB.Close()

			_, err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName))
			if err != nil {
				return nil, fmt.Errorf("failed to create database: %w", err)
			}

			err = db.Ping()
			if err != nil {
				return nil, fmt.Errorf("failed to connect to newly created database: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sqlx.DB) error {

	migrationSQL, err := ioutil.ReadFile("internal/adapters/output/db/migrations/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	statements := strings.Split(string(migrationSQL), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing migration statement: %w\nStatement: %s", err, statement)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
