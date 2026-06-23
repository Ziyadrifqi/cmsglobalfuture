package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Tambahkan ini
	log.Println("DB_HOST =", cfg.DBHost)
	log.Println("DB_PORT =", cfg.DBPort)
	log.Println("DB_USER =", cfg.DBUser)
	log.Println("DB_PASSWORD =", cfg.DBPassword)
	log.Println("DB_NAME =", cfg.DBName)

	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}

	log.Println("✓ Database terhubung")
	return db
}
