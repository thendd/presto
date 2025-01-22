package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Connection *pgx.Conn

func Connect() {
	log.Println("Started connecting to the database")
	connectionString := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB")

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
