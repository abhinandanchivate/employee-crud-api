package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/abhinandanchivate/employee-crud-api/internal/config"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
)

type Database struct {
	DB *gorm.DB
}

var dbInstance *Database

func Init() (*Database, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	cfg := config.Get()

	// Create DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
		cfg.Database.ParseTime,
		cfg.Database.Loc,
	)

	// Configure GORM logger based on environment
	gormLogger := logger.Default
	if cfg.App.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.Employee{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Database connected successfully: %s", cfg.Database.Name)

	dbInstance = &Database{DB: db}
	return dbInstance, nil
}

func Get() *gorm.DB {
	if dbInstance == nil {
		_, err := Init()
		if err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
	}
	return dbInstance.DB
}

func Close() error {
	if dbInstance != nil && dbInstance.DB != nil {
		sqlDB, err := dbInstance.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
