package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"cryptorate-service/internal/bot"
	_ "github.com/lib/pq"
)

func main() {
	//Подключение к БД
	connStr := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	//Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	bot, err := bot.NewBot(token, db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started...")
	bot.Start()
}
