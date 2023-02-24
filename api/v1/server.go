package v1

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ishanshre/go-auth-api/api/v1/models"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type Storage interface {
	UserSignUp(*models.RegisterUser) error
	UserLogin(string) (*models.UserNhash, error)
	AdminUpdateUserById(int, *models.AdminUpdateUser) error
	AdminDeleteUserById(int) error
	GetUsers() ([]*models.User, error)
	GetUsersById(int) (*models.User, error)
	DeleteUserById(int) error
	UpdateUserById(int, *models.UpdateUser) error
	UpdateLastLogin(int) error
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

func (s *PostgresStore) UserSignUp(user *models.RegisterUser) error {
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

func (s *PostgresStore) GetUsersById(id int) (*models.User, error) {
	query := `
		SELECT * FROM users 
		WHERE id = $1
	`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanUsers(rows)
	}
	return nil, fmt.Errorf("account with id %v not found", id)
}

func (s *PostgresStore) DeleteUserById(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	s.db.Exec("COMMIT")
	rows, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows_affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return fmt.Errorf("id %v does not exists", id)
	}
	return nil
}
func (s *PostgresStore) UpdateUserById(id int, user *models.UpdateUser) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3
		WHERE id =$1
	`
	s.db.Exec("COMMIT")
	rows, err := s.db.Exec(query, id, user.FistName, user.LastName)
	if err != nil {
		return err
	}
	rows_affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return fmt.Errorf("id %v does not exists", id)
	}
	return nil
}

func (s *PostgresStore) UserLogin(username string) (*models.UserNhash, error) {
	query := `
		SELECT id, password FROM users
		WHERE username = $1
	`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanUser1(rows)
	}
	return nil, fmt.Errorf("username: %v not found", username)
}

func (s *PostgresStore) UpdateLastLogin(id int) error {
	query := `
		UPDATE users 
		SET last_login = $2
		WHERE id = $1
	`
	s.db.Exec("COMMIT")
	_, err := s.db.Query(query, id, time.Now())
	return err
}

func (s *PostgresStore) AdminUpdateUserById(id int, user *models.AdminUpdateUser) error {
	query := `
		UDPATE users 
		SET first_name = $2, last_name= $3, is_admin = $4, is_active =$5, updated_at = $6
		WHERE id = $1
	`
	s.db.Exec("COMMIT")
	rows, err := s.db.Exec(query, id, user.FirstName, user.LastName, user.IsAdmin, user.IsActive, time.Now())
	if err != nil {
		return err
	}
	rows_affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return fmt.Errorf("id %v does not exists", id)
	}
	return nil
}

func (s *PostgresStore) AdminDeleteUserById(id int) error {
	return nil
}
