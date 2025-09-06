package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func DbConnection() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, dbName,
	)
	// Connect to DB
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}
	if err := migrateUp(db, dsn); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logrus.Info("✅ DB connection successful")
	return db, nil
}

// findMigrationsFolderRoot searches upward until migrations folder is found
func findMigrationsFolderRoot() string {
	wd, _ := os.Getwd()
	relPath := "database/migrations"

	for {
		currentPath := filepath.Join(wd, relPath)
		if fi, err := os.Stat(currentPath); err == nil && fi.IsDir() {
			return currentPath
		}
		newDir := filepath.Dir(wd)
		if newDir == "/" || newDir == wd {
			break
		}
		wd = newDir
	}
	return ""
}

// migrateUp applies pending migrations
// migrateUp applies pending migrations
func migrateUp(db *sqlx.DB, dsn string) error {
	sqlDB := db.DB
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("create DB driver failed: %w", err)
	}

	path := findMigrationsFolderRoot()
	if path == "" {
		return fmt.Errorf("migrations folder not found")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration instance creation failed: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logrus.Info("ℹ️ No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	logrus.Info("✅ Database migrated successfully")
	return nil
}
