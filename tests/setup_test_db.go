package tests

import (
	"dklautomationgo/database"
	"fmt"
	"os"

	"gorm.io/gorm"
)

// SetupTestDB initializes a test database connection
func SetupTestDB() (*gorm.DB, error) {
	// Use a test database configuration
	dbConfig := database.NewConfig()
	dbConfig.DBName = "dklautomationgo_test"

	// Override with environment variables if provided
	if os.Getenv("TEST_DB_HOST") != "" {
		dbConfig.Host = os.Getenv("TEST_DB_HOST")
	}
	if os.Getenv("TEST_DB_PORT") != "" {
		dbConfig.Port = os.Getenv("TEST_DB_PORT")
	}
	if os.Getenv("TEST_DB_USER") != "" {
		dbConfig.User = os.Getenv("TEST_DB_USER")
	}
	if os.Getenv("TEST_DB_PASSWORD") != "" {
		dbConfig.Password = os.Getenv("TEST_DB_PASSWORD")
	}
	if os.Getenv("TEST_DB_NAME") != "" {
		dbConfig.DBName = os.Getenv("TEST_DB_NAME")
	}

	// Connect to the test database
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Migrate schema
	err = database.AutoMigrate(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate test database: %w", err)
	}

	return db, nil
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(db *gorm.DB) error {
	// Get the underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Close the connection
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// CleanupTestData removes all data from test tables
func CleanupTestData(db *gorm.DB) error {
	// List of tables to clean in reverse order of dependencies
	tables := []string{
		"refresh_tokens",
		"users",
		"aanmeldingen",
		"contact_formulieren",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
