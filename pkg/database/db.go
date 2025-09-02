package database

import (
	"context"
	"fmt"
	"time"

	"booking/pkg/config"
	"booking/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDB(config *config.DatabaseConfig, log logger.Logger) *sqlx.DB {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
	db, err := sqlx.Open("pgx", connString)
	if err != nil {
		log.Fatalf(err, "error opening database connection")
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf(err, "error pinging database")
	}

	log.Infof("Database pool initialized successfully - MaxOpenConns: %d, MaxIdleConns: %d",
		config.MaxOpenConns, config.MaxIdleConns)

	return db
}
