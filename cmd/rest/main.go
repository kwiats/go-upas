package main

import (
	"log"

	"github.com/kwiats/go-upas/server"
	"github.com/kwiats/go-upas/storage"
)

func main() {
	store := storage.CreateConnectionDatabase()
	log.Print(store)
	log.Print("User Profile Auth System")

	server := server.RunAPIServer(":3000", store)
	server.Run()
}
