package migrations

import (
	"fmt"
	"log"
	// "path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(dbURL string, migrationsPath string) error {
	fmt.Println("Running migrations from:", migrationsPath)
	fmt.Println("Connecting to DB:", dbURL)
	// absPath := filepath.Join("notes-service", "internal", "migrations")
	// path := filepath.ToSlash(filepath.Clean(absPath))
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	log.Printf("Applying migrations from: file:///%s", migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
