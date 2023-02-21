package v1

import (
	"database/sql"
	"log"
	"os"

	"github.com/ishanshre/go-auth-api/api/v1/models"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type Storage interface {
	CreateUser(*models.User) error
	GetUsers() ([]*models.User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalf("Error in loading environment files: %v", err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createUserTable()
}

func (s *PostgresStore) createUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(500) NOT NULL,
			is_admin BOOLEAN DEFAULT 'f',
			is_active BOOLEAN DEFAULT 't',
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			last_login TIMESTAMPTZ
		)
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			username,
			email,
			password,
			created_at,
			updated_at,
			last_login
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLogin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetUsers() ([]*models.User, error) {
	query := `SELECT * FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*models.User{}
	for rows.Next() {
		user, err := scanUsers(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanUsers(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)
	return user, err
}
