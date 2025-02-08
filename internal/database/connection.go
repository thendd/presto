package database

import (
	"log"
	"presto/internal/database/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Connect() {
	log.Println("Started connecting to the database")

	conn, err := gorm.Open(postgres.Open(config.POSTGRESQL_CONNECTION_STRING))
	if err != nil {
		log.Println("Could not connect to the database: ", err)
	}

	Connection = conn

	Connection.AutoMigrate(&Guild{})

	log.Println("Connected to the database successfully")
}
