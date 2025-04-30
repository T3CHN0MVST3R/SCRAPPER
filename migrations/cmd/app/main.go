package main

import (
	"database/sql"
	"fmt"
	"log"

	"scrapper-migrator/config"

	_ "github.com/lib/pq"
	"github.com/vnlozan/goose/v3"
)

func main() {
	cfg := config.NewConfig()

	fmt.Println("Database DSN:", cfg.DB.DSN)
	fmt.Println("Migrations directory:", cfg.MigrateDIR)

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	if err = goose.Up(db, cfg.MigrateDIR); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
