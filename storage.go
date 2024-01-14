package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	DeleteUser(string) error
	UpdateUser(*User) error
	GetUserByID(string) (*User, error)
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
	dbs.createProfilTable()
	dbs.createUserTable()
	return nil
}

func (dbs *DataBaseStore) createProfilTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS profils (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
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
		updated_at TIMESTAMP WITH TIME ZONE,
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255),
		email VARCHAR(255) UNIQUE,
		is_admin BOOLEAN,
		is_active BOOLEAN,
		profil_id UUID REFERENCES profils(id)
	)`
	log.Println(query)
	_, err := dbs.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DataBaseStore) CreateUser(user *User) error {
	transaction, err := dbs.db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}
	profileQuery := `INSERT INTO profils (created_at, updated_at, image_profile, first_name, last_name)
                 VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var profilID string
	err = transaction.QueryRow(
		profileQuery,
		user.CreatedAt,
		user.UpdatedAt,
		user.Profil.ImageProfile,
		user.Profil.FirstName,
		user.Profil.LastName,
	).Scan(&profilID)

	if err != nil {
		transaction.Rollback()
		log.Printf("Failed to insert profile: %v", err)
		return err
	}

	userQuery := `INSERT INTO users ( created_at, updated_at, username, password, email, is_admin, is_active, profil_id)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var userID string
	err = transaction.QueryRow(
		userQuery,
		user.CreatedAt,
		user.UpdatedAt,
		user.Username,
		user.Password,
		user.Email,
		user.IsAdmin,
		user.IsActive,
		profilID,
	).Scan(&userID)

	if err != nil {
		transaction.Rollback()
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	if err := transaction.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	log.Printf("Created user(UUID=%v) and profile(UUID=%v) object.", userID, profilID)
	return nil
}
func (dbs *DataBaseStore) UpdateUser(*User) error {
	return nil
}
func (dbs *DataBaseStore) DeleteUser(id string) error {
	return nil
}
func (dbs *DataBaseStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}
