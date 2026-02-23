package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

func MustConnectToDb(cfg Config) *sql.DB {
	connstr := fmt.Sprintf(`host=%s port=%s username=%s password=%s dbname=%s sslmode=disable`,
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
