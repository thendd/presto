package database

import (
	"presto/internal/database/config"
	"presto/internal/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Connect() {
	log.Info("Started connecting to the database")

	conn, err := gorm.Open(postgres.Open(config.POSTGRESQL_CONNECTION_STRING))
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	Connection = conn

	Connection.AutoMigrate(&Guild{})

	log.Info("Connected to the database successfully")
}
