package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/CodingFervor/carbon-emission-management/internal/config"
	"github.com/CodingFervor/carbon-emission-management/pkg/logger"
)

var DB *sql.DB

func Connect(cfg config.DatabaseConfig) error {
	var err error
	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(2 * time.Minute)
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	logger.Info("database connected", "host", cfg.Host, "db", cfg.DBName)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		logger.Info("database connection closed")
	}
}
