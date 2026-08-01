package db

import (
	"context"
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GDB *gorm.DB

func Connection() error {
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := os.Getenv("DB_PASSWORD")
	dbname := strings.TrimSpace(os.Getenv("DB_NAME"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		return errors.New("faltan variables obligatorias de PostgreSQL")
	}
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return errors.New("DB_PORT no es valido")
	}

	sslMode, err := databaseSSLMode(host)
	if err != nil {
		return err
	}
	connectionURL := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + dbname,
	}
	query := connectionURL.Query()
	query.Set("sslmode", sslMode)
	query.Set("connect_timeout", "5")
	connectionURL.RawQuery = query.Encode()

	logMode := logger.Warn
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production") {
		logMode = logger.Info
	}
	GDB, err = gorm.Open(postgres.Open(connectionURL.String()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logMode),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	sqlDB, err := GDB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}
	log.Println("BD conectada")
	return nil
}

func databaseSSLMode(host string) (string, error) {
	configured := strings.ToLower(strings.TrimSpace(os.Getenv("DB_SSLMODE")))
	if configured == "" {
		if isLoopbackHost(host) {
			return "disable", nil
		}
		return "", errors.New("DB_SSLMODE es obligatorio para una base de datos remota")
	}
	allowed := map[string]bool{
		"disable": true, "require": true, "verify-ca": true, "verify-full": true,
	}
	if !allowed[configured] {
		return "", errors.New("DB_SSLMODE no es valido")
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production") && !isLoopbackHost(host) && configured != "verify-ca" && configured != "verify-full" {
		return "", errors.New("una base de datos remota en produccion requiere DB_SSLMODE=verify-ca o verify-full")
	}
	return configured, nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return ip != nil && ip.IsLoopback()
}
