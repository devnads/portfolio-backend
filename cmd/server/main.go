package main

import (
	"context"
	"log"
	"monad-indexer/internal/db"
	"monad-indexer/internal/routes"
	"net/http"
)


func main() {
	db.InitDB()
	defer db.Conn.Close(context.Background())

	db.Migrate()

	r := routes.SetupRoutes()

	log.Println("Server running on port :8080")
	http.ListenAndServe(":8080",r)
}





