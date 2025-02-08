package config

import (
	"os"
)

var (
	POSTGRESQL_CONNECTION_STRING string
)

func Load() {
	POSTGRESQL_CONNECTION_STRING = "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB")
}
