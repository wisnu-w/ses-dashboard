package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var (
		dbHost     = flag.String("host", "localhost", "Database host")
		dbPort     = flag.String("port", "5432", "Database port")
		dbUser     = flag.String("user", "ses_user", "Database user")
		dbPassword = flag.String("password", "ses_password", "Database password")
		dbName     = flag.String("dbname", "ses_monitoring", "Database name")
		action     = flag.String("action", "up", "Migration action: up, down")
	)
	flag.Parse()

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migrate instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infrastructure/database/migration",
		"postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Execute action
	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations up:", err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations down:", err)
		}
		fmt.Println("Migrations rolled back successfully")
	default:
		log.Fatal("Unknown action:", *action)
	}
}