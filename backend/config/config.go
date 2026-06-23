package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	SessionSecret  string
	JWTSecret      string
	AllowedOrigins string

	// ── SMTP (untuk kirim email akun baru) ──────────────────────────────────
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFromName string
	AppBaseURL   string // dipakai untuk link "Masuk ke CMS" di email
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env tidak ditemukan, pakai env sistem")
	}

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		smtpPort = 587
	}

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "cms_db"),
		SessionSecret:  getEnv("SESSION_SECRET", "rahasia-session"),
		JWTSecret:      getEnv("JWT_SECRET", "rahasia-jwt"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "Yayasan CMS"),
		AppBaseURL:   getEnv("APP_BASE_URL", "http://localhost:8080"),
	}

	if cfg.AppEnv == "production" && cfg.SessionSecret == "rahasia-session" {
		log.Println("⚠ PERINGATAN: SESSION_SECRET masih memakai nilai default di mode production.")
		log.Println("⚠ Set SESSION_SECRET ke nilai acak yang konsisten di semua instance.")
	}

	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		log.Println("⚠ SMTP_USER / SMTP_PASSWORD belum diisi di .env — email notifikasi akun baru TIDAK akan terkirim.")
		log.Println("⚠ Gunakan Gmail App Password: https://myaccount.google.com/apppasswords")
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
