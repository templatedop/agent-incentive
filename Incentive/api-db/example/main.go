package main

import (
	"log"

	dblib "gitlab.cept.gov.in/it-2.0-common/n-api-db"
)

func main() {
	// Step 1: Initialize the DB configuration
	dbConfig := dblib.DBConfig{
		DBUsername:        "db_username",
		DBPassword:        "db_password",
		DBHost:            "your_host_name",
		DBPort:            "5432",
		DBDatabase:        "your_database_name",
		Schema:            "your_db_schema",
		MaxConns:          10,
		MinConns:          2,
		MaxConnLifetime:   30, // In minutes
		MaxConnIdleTime:   10, // In minutes
		HealthCheckPeriod: 5,  // In minutes
		AppName:           "your_app_name",
		SSLMode:           "disable",
	}
	// Step 2: Prepare the DB configuration (validate and set defaults)
	preparedConfig := dblib.NewDefaultDbFactory().NewPreparedDBConfig(dbConfig)

	// Step 3: Establish the database connection
	db, err := dblib.NewDefaultDbFactory().CreateConnection(preparedConfig, nil, nil)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}
	defer db.Close() // Ensure the connection is closed when the program exits

	if db.Ping() == nil {
		log.Printf("Database connection established successfully")

	} else {

		log.Printf("Error pinging the database connection")
	}
}
