package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dbURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	// Connection pool (rất quan trọng)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	log.Println("PostgreSQL connected")
	return db
}
