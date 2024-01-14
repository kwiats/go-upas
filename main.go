package main

import (
	"log"
)

func CreateConnectionDatabase() *DataBaseStore {
	log.Print("Connection to Postgres database...")
	store, err := CreateDBConnection()
	if err != nil {
		log.Fatalln("Connection failed.", err)
	}
	if err := store.Init(); err != nil{
		log.Printf("Failed to create tables: %v", err)

	}
	log.Print("Connected to Postgres database.")
	return store
}

func main() {
	store := CreateConnectionDatabase()
	log.Print(store)
	log.Print("User Profile Auth System")

	server := runAPIServer(":3000", &store)
	server.Run()
}
