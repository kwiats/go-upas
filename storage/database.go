package storage

import (
	"database/sql"
	"github.com/kwiats/go-upas/models"
	"github.com/kwiats/go-upas/utils"
	"log"
)

type DataBaseStore struct {
	db *sql.DB
}

// SignInUser implements Storage.
func (*DataBaseStore) SignInUser(*models.UserLogin) (*utils.TokenResponse, error) {
	panic("unimplemented")
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

func CreateConnectionDatabase() *DataBaseStore {
	log.Print("Connection to Postgres database...")
	store, err := CreateDBConnection()
	if err != nil {
		log.Fatalln("Connection failed.", err)
	}
	if err := store.Init(); err != nil {
		log.Printf("Failed to create tables: %v", err)

	}
	log.Print("Connected to Postgres database.")
	return store
}
