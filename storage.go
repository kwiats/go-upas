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
		id UUID PRIMARY KEY,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		image_profile TEXT,
		first_name VARCHAR(255),
		last_name VARCHAR(255)
	)`
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
		id UUID PRIMARY KEY,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE,
		username VARCHAR(255),
		password VARCHAR(255),
		email VARCHAR(255),
		is_admin BOOLEAN,
		is_active BOOLEAN,
		profil_id UUID REFERENCES profils(id)
	)`
	_, err := dbs.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DataBaseStore) CreateUser(*User) error {
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
