package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
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
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

func Login(user models.Login_Request, db *sql.DB) (models.User, error) {
	query := `SELECT * FROM users WHERE login = $1`
	ans := db.QueryRow(query, user.Login)
	return ans
}

func Registration(user models.Registration_Request) *sql.DB {

}
