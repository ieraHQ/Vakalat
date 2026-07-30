package main

import (
	"log"
	"os"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: goose [command] [args]")
	}

	command := os.Args[1]
	args := os.Args[2:]

	dbURL := "postgres://postgres:postgres@localhost:5432/vakalat?sslmode=disable"
	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close DB: %v", err)
		}
	}()

	dir := "./database/migrations"
	if err := goose.Run(command, db, dir, args...); err != nil {
		log.Fatalf("Failed to run goose: %v", err)
	}
}