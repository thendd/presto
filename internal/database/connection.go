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

	conn, err := pgx.Connect(context.Background(), config.POSTGRESQL_CONNECTION_STRING)
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	Connection = conn

	rawSchema, _ := os.ReadFile("./internal/database/schema.sql")

	_, err = Connection.Exec(context.Background(), string(rawSchema))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database successfully")
}
