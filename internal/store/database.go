package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "DB_PORT", "DB_SSL_MODE"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	ssl := os.Getenv("DB_SSL_MODE")
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host,
			user,
			pass,
			name,
			port,
			ssl,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("Error opening db: %w", err)
	}

	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := waitForDatabase(db); err != nil {
		return nil, err
	}

	fmt.Println("Connected to database...")
	return db, nil
}

func waitForDatabase(db *sql.DB) error {
	const (
		maxAttempts    = 10
		initialBackoff = 500 * time.Millisecond
		maxBackoff     = 5 * time.Second
	)

	backoff := initialBackoff
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := db.Ping()
		if err == nil {
			return nil
		}

		lastErr = err
		fmt.Printf("Waiting for database (attempt %d/%d): %v\n", attempt, maxAttempts, err)
		time.Sleep(backoff)

		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}

	return fmt.Errorf("Error pinging db after %d attempts: %w", maxAttempts, lastErr)
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: set dialect: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: goose up: %w", err)
	}
	return nil
}
