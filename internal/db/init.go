package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Conn *pgxpool.Pool

func InitDB() {
	err := godotenv.Load("/home/ubuntu/.env")
	if err != nil {
		log.Println("Erreur reading .env file.")
		return
	}

	url := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("DB pool connection failed:", err)
	}

	Conn = pool
	fmt.Println("✅ DB Pool connecté")
}
