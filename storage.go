package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*UserCreate) error
	DeleteUser(string) error
	UpdateUser(*UserDTO) error
	GetUsers() ([]*UserDTO, error)
	GetUserByID(string) (*UserDTO, error)
}

type DataBaseStore struct {
	db *sql.DB
}

func CreateDBConnection() (*DataBaseStore, error) {

	connStr := "postgres://postgres:123456789@localhost:5432/postgres?sslmode=disable" // delete hardcode fields
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DataBaseStore{
		db: db,
	}, nil
}

func (dbs *DataBaseStore) Init() error {
	err := dbs.createProfilTable()
	if err != nil {
		return err
	}
	err = dbs.createUserTable()
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DataBaseStore) createProfilTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS profils (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		image_profile TEXT,
		first_name VARCHAR(255),
		last_name VARCHAR(255)
	)`
	log.Println(query)
	_, err := dbs.db.Exec(query)
	if err != nil {
		log.Printf("Failed to create Profil table: %v", err)
		return err
	}
	return nil
}

func (dbs *DataBaseStore) createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255),
		email VARCHAR(255) UNIQUE,
		is_admin BOOLEAN DEFAULT FALSE,
		is_active BOOLEAN DEFAULT FALSE,
		profil_id UUID REFERENCES profils(id)
	)`
	log.Println(query)
	_, err := dbs.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DataBaseStore) CreateUser(user *UserCreate) error {
	transaction, err := dbs.db.Begin()
	defer func() {
		if err != nil {
			if rErr := transaction.Rollback(); rErr != nil {
				log.Printf("Error rolling back transaction: %v", rErr)
			}
			log.Printf("Transaction rolled back: %v", err)
		}
	}()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	profilID, err := InsertProfile(transaction, user)

	if err != nil {
		log.Printf("Failed to insert profile: %v", err)
		return err
	}
	userID, err := InsertUser(transaction, user, profilID)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	if err := transaction.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	log.Printf("Created user(UUID=%+v) and profile(UUID=%+v) object.", userID, profilID)
	return nil
}

func InsertUser(transaction *sql.Tx, user *UserCreate, profilID string) (string, error) {
	userQuery := `INSERT INTO users ( username, password, email,  profil_id)
				VALUES ($1, $2, $3, $4) RETURNING id`

	var userID string
	err := transaction.QueryRow(
		userQuery,
		user.Username,
		user.Password,
		user.Email,
		profilID,
	).Scan(&userID)
	return userID, err
}

func InsertProfile(transaction *sql.Tx, user *UserCreate) (string, error) {
	profileQuery := `INSERT INTO profils (image_profile, first_name, last_name)
                 VALUES ($1, $2, $3) RETURNING id`

	var profilID string
	err := transaction.QueryRow(
		profileQuery,
		user.ImageProfile,
		user.FirstName,
		user.LastName,
	).Scan(&profilID)
	return profilID, err
}

func (dbs *DataBaseStore) GetUsers() ([]*UserDTO, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := dbs.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := []*UserDTO{}
	for rows.Next() {
		user := new(UserDTO)
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbs *DataBaseStore) UpdateUser(*UserDTO) error {
	return nil
}

func (dbs *DataBaseStore) DeleteUser(id string) error {
	return nil
}
func (dbs *DataBaseStore) GetUserByID(id string) (*UserDTO, error) {
	return &UserDTO{}, nil
}
