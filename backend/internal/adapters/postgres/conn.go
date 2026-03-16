package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) GetByLogin(login string) (models.User, error) {
	query := `
		SELECT id, username, login, password_hash, created_at
		FROM users
		WHERE login = $1
	`

	var ans models.User
	err := r.DB.QueryRow(query, login).Scan(
		&ans.ID,
		&ans.Username,
		&ans.Login,
		&ans.Password_Hash,
		&ans.Created_at,
	)
	return ans, err
}

func (r *UserRepo) Create(user models.User) error {
	query := `INSERT INTO users (id, username, login, password_hash, created_at, telegram_id) 
              VALUES ($1, $2, $3, $4, NOW(), $5)`
	_, err := r.DB.Exec(
		query,
		user.ID,
		user.Username,
		user.Login,
		user.Password_Hash,
		user.Telegram_id,
	)

	var myerror *pq.Error
	if errors.As(err, &myerror) {
		if myerror.Code == "23505" {
			return errors.New("user already exists")
		}
	}
	return err
}

func (r *UserRepo) GetByTelegramID(id int64) (models.User, error) {
	query := `SELECT id, username, login, password_hash, created_at, telegram_id 
			FROM users 
			WHERE telegram_id = $1`

	var user models.User

	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Login, &user.Password_Hash, &user.Created_at, &user.Telegram_id)
	return user, err
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
