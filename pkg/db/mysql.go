package db

import (
	"blog-api/config"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func ConnectDB(cfg *config.Config) *sql.DB {
	// Connect to MySQL server (without specifying the DB name yet)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the MySQL server:", err)
	}

	// Ensure the connection is established
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the MySQL server:", err)
	}

	// Create the database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.DBName))
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	// Now connect to the specific database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Ping the database again to ensure the connection is valid
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Println("Connected to the database successfully.")
	return db
}

func InitializeDB(db *sql.DB, cfg *config.Config) error {

	// Read SQL script for creating tables
	sqlFile, err := ioutil.ReadFile("scripts/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL script: %v", err)
	}

	// Split SQL statements (assuming statements are separated by semicolons)
	statements := strings.Split(string(sqlFile), ";")

	for i, stmt := range statements {
		if stmt != "" {
			_, err = db.Exec(stmt)
			if err != nil {
				log.Printf("Error executing statement %d: %v", i+1, err)
				return fmt.Errorf("failed to execute SQL statement %d: %v", i+1, err)
			}
		}
	}

	log.Println("Tables initialized successfully.")
	return nil
}
