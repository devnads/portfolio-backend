package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func InitDB() {
	// err := godotenv.Load("../../.env") LOCAL TEST
	err := godotenv.Load("/home/ubuntu/.env")
	if err != nil {
		log.Println("Erreur reading .env file.")
		return
	}

	url := os.Getenv("DB_URL")
	conn, err := pgx.Connect(context.Background(),url )
	if err != nil {
		log.Fatal("DB connexion failed", err)
	}
	Conn = conn
	fmt.Println("DB Connecté")
}