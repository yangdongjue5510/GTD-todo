package db

import (
	"fmt"
	"time"
	"yangdongju/gtd_todo/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnectionPool(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection pool 설정
	db.SetMaxOpenConns(cfg.DBMaxOpenConnection)
	db.SetMaxIdleConns(cfg.DBMaxIdleConnection)
	db.SetConnMaxLifetime(time.Duration(cfg.DBConnectionMaxLifeTime) * time.Minute)

	// 실제 연결 확인
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
