package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"arczed/internal/server/configs"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	MainConnect() *gorm.DB // เชื่อมต่อกับฐานข้อมูลหลัก
	Close() error          // ปิดการเชื่อมต่อ
}

type service struct {
	db *gorm.DB
}

var (
	dbInstance *service
)

// New initializes a new GORM service.
func New(config *configs.Config) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	var (
		database = config.DbDatabase
		password = config.DbPassword
		username = config.DbUsername
		port     = config.DbPort
		host     = config.DbHost
	)

	// Connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	// Open the connection using GORM with PostgreSQL driver
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // ตั้งค่า log ให้เงียบหรือ debug ได้
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB() // Retrieve *sql.DB from GORM for connection pool configuration
	if err != nil {
		log.Fatal(err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Save the instance for reuse
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to retrieve sql.DB: %v", err)
		log.Fatalf("failed to retrieve sql.DB: %v", err)
		return stats
	}

	// Ping the database
	if err := sqlDB.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()

	return stats
}

// MainConnect returns the GORM database instance.
func (s *service) MainConnect() *gorm.DB {
	return s.db
}

// Close closes the database connection.
func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve sql.DB: %w", err)
	}
	log.Printf("Disconnected from database: %s", "Main DataBase")
	return sqlDB.Close()
}
