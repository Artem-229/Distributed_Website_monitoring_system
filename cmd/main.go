package main

import (
	"Distributed_Website_monitoring_system/internal/app"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	envinf := app.MustGetFromEnv()

	connstr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		envinf.DB_USERNAME, envinf.DB_PASSWORD, envinf.DB_PORT, envinf.DB_NAME)

	m, err := migrate.New(
		"file://migrations",
		connstr,
	)

	if err != nil {
		panic(err)
	}

	m.Up()

	app.Run(app.MustGetFromEnv())
}
