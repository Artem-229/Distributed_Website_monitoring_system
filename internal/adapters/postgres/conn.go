package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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
	var ans models.User
	err := db.QueryRow(query, user.Login).Scan(
		&ans.ID,
		&ans.Username,
		&ans.Login,
		&ans.Password_Hash,
		&ans.Created_at,
	)
	return ans, err
}

func Registration(user models.Registration_Request, db *sql.DB) error {
	query := `INSERT INTO users (id, username, login, password_hash, created_at) 
              VALUES ($1, $2, $3, $4, NOW())`
	newID := uuid.New()
	_, err := db.Exec(query, newID, user.Username, user.Login, user.Password)
	return err
}
