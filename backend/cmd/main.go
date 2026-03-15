package main

import (
	"Distributed_Website_monitoring_system/internal/adapters/postgres"
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/controller"
	"Distributed_Website_monitoring_system/internal/handlers"
	"Distributed_Website_monitoring_system/internal/kafka/producer"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	envinf := app.MustGetFromEnv()

	secret := envinf.JWT_SECRET

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		envinf.DB_USERNAME, envinf.DB_PASSWORD, envinf.DB_HOST, envinf.DB_PORT, envinf.DB_NAME)
	m, err := migrate.New("file:///app/migrations", connstr)

	if err != nil {
		panic(err)
	}

	fmt.Println("START MIGRATION")

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("MIGRATION ERROR:", err)
	}

	fmt.Println("MIGRATION DONE")

	cfg := postgres.Config{
		Host:     envinf.DB_HOST,
		Port:     envinf.DB_PORT,
		Username: envinf.DB_USERNAME,
		Password: envinf.DB_PASSWORD,
		Dbname:   envinf.DB_NAME,
	}

	conn := postgres.MustConnectToDb(cfg)

	userrepo := &postgres.UserRepo{
		DB: conn,
	}

	monitorrepo := &postgres.MonitorRepo{
		DB: conn,
	}

	checksrepo := &postgres.ChecksRepo{
		DB: conn,
	}

	authHandler := handlers.NewAuthHandler(userrepo, secret)
	monitorHandler := handlers.NewMonitorHandler(monitorrepo)
	checkHandler := handlers.NewCheckHandler(checksrepo)
	healthHandler := &handlers.HealthHandler{}

	producer, err := producer.NewProducer(envinf.KAFKA_ADDRESS)
	if err != nil {
		fmt.Println("Kafka error", err)
	}

	go func() {
		for {
			arr, err := monitorrepo.GetAllMonitors()
			fmt.Println("WORKER: monitors count:", len(arr), "err:", err)
			if err == nil {
				for _, k := range arr {
					_, _, err := app.CheckPing(k, checksrepo, producer)
					fmt.Println("WORKER: checked", k.Url, "err:", err)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	contrl := controller.SetupRoutes(authHandler, monitorHandler, checkHandler, healthHandler, secret)

	contrl.Listen(":8080")
}
