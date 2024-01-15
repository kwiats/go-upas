package storage

import (
	"log"

	_ "github.com/lib/pq"
)

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
		profil_id UUID REFERENCES profils(id) ON DELETE CASCADE
	)`
	log.Println(query)
	_, err := dbs.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
