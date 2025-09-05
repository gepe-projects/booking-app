package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(db *sql.DB, dsn string) {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatalf("failed to open migrations: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange { // kalo errornya errNoChange, berarti gada perubahan
		log.Fatalf("failed to run migrate up: %v", err)
	}

	fmt.Println("✅ Database migrated successfully")
}
