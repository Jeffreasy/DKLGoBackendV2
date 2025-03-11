package database

import (
	"dklautomationgo/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config bevat de database configuratie
type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

// NewConfig maakt een nieuwe database configuratie op basis van omgevingsvariabelen
func NewConfig() *Config {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable" // Standaard voor lokale ontwikkeling
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "dklautomationgo"
	}

	config := &Config{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbName,
		Port:     port,
		SSLMode:  sslMode,
	}

	log.Printf("[Database] Configuration: Host=%s, Port=%s, User=%s, Database=%s, SSLMode=%s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode)

	return config
}

// DSN geeft de database connection string terug
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Amsterdam",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}

// NewConnection maakt een nieuwe database verbinding
func NewConnection(config *Config) (*gorm.DB, error) {
	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	dsn := config.DSN()
	log.Printf("[Database] Connecting to PostgreSQL with DSN: %s",
		// Verberg het wachtwoord in de logs
		fmt.Sprintf(
			"host=%s user=%s password=*** dbname=%s port=%s sslmode=%s TimeZone=Europe/Amsterdam",
			config.Host, config.User, config.DBName, config.Port, config.SSLMode,
		),
	)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test de verbinding
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configureer de connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("[Database] Successfully connected to PostgreSQL database at %s:%s/%s",
		config.Host, config.Port, config.DBName)

	return db, nil
}

// AutoMigrate voert automatische migraties uit voor de opgegeven modellen
func AutoMigrate(db *gorm.DB) error {
	log.Println("[Database] Running auto migrations...")

	// Voeg hier je modellen toe
	err := db.AutoMigrate(
		&models.ContactFormulier{},
		&models.Aanmelding{},
		&models.User{},
		&models.RefreshToken{},
	)

	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("[Database] Auto migrations completed successfully")
	return nil
}
