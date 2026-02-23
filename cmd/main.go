package main

import (
	"Distributed_Website_monitoring_system/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	app.Run(app.MustGetFromEnv())
}
