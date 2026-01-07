package db

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dbURL string) *sqlx.DB {
	// Log connection info (without password) for debugging
	if dbURL != "" {
		// Mask password in log
		maskedURL := dbURL
		if idx := strings.Index(dbURL, "@"); idx > 0 {
			if passIdx := strings.LastIndex(dbURL[:idx], ":"); passIdx > 0 {
				maskedURL = dbURL[:passIdx+1] + "***" + dbURL[idx:]
			}
		}
		log.Printf("Connecting to database: %s", maskedURL)
	} else {
		log.Println("WARNING: DATABASE_URL is empty!")
	}
	
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
