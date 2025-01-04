package database

import (
	"context"
	"log"
	"os"

	"presto/internal/database/config"

	"github.com/jackc/pgx/v5"
)

var Connection *pgx.Conn

func Connect() {
	log.Println("Started connecting to the database")
	connectionString := config.POSTGRESQL_CONNECTION_STRING

	if os.Getenv("ENVIRONMENT") == "development" {
		connectionString = config.DEVELOPMENT_POSTGRESQL_CONNECTION_STRING
	}

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	Connection = conn

	rawSchema, _ := os.ReadFile("./internal/database/schema.sql")

	_, err = conn.Exec(context.Background(), string(rawSchema))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database successfully")
}
