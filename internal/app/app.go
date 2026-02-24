package app

import (
	"Distributed_Website_monitoring_system/internal/adapters/postgres"
	"Distributed_Website_monitoring_system/internal/controller"
)

func Run(cfg Config) {
	conn := postgres.MustConnectToDb(
		postgres.Config{
			Host:     cfg.DB_HOST,
			Port:     cfg.DB_PORT,
			Username: cfg.DB_USERNAME,
			Password: cfg.DB_PASSWORD,
			Dbname:   cfg.DB_NAME,
		},
	)

	_ = conn

	contrl := controller.SetupRoutes()

	contrl.Listen(":8080")

}
