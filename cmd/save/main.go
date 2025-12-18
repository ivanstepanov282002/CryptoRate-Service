package main

import (
	"context"
	"cryptorate-service/internal/api"
	"cryptorate-service/internal/models"
	"cryptorate-service/internal/repository"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

// Запус автоматической выгрузки по API курса валют с промежутком времени interval секунды
func main() {
	interval := flag.Int("interval", 0, "Update interval in MINUTES (0 = run once)")
	flag.Parse()

	// Подключение к БД
	// Добавлен fallback на значения по умолчанию
	connStr := "host=127.0.0.1 port=5432 user=crypto_user password=secure_password_123 dbname=crypto_db sslmode=disable"

	// Пробуем получить из .env, если не получилось - используем значения выше
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		connStr = fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"))
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}
	fmt.Println("✅ Connected to database")

	repo := repository.NewRepository(db)
	client := api.NewCoinGeckoClient()

	if *interval == 0 {
		// Одноразовый запуск
		fmt.Println("🚀 One-time rates update")
		updateRates(client, repo)
	} else {
		fmt.Printf("🚀 Worker started. Fetching rates every %d minutes...\n", *interval)
		fmt.Println("Press Ctrl+C to stop")

		//Добавлен graceful shutdown для мягкой остановки
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		//Создает канал, который будет отсылать текущее время с периодичностью interal
		ticker := time.NewTicker(time.Duration(*interval) * time.Minute)
		defer ticker.Stop()

		// Первый запуск сразу
		updateRates(client, repo)

		for {
			select {
			case <-ticker.C:
				updateRates(client, repo)
			case <-ctx.Done():
				fmt.Println("\n👋 Stopping worker...")
				return
			}
		}
	}
}

func updateRates(client *api.CoinGeckoClient, repo *repository.Repository) {
	//Добавлен timestamp в логи
	currentTime := time.Now().Format("15:04")
	fmt.Printf("\n⏰ [%s] Fetching rates for 7 currencies...\n", currentTime)

	coinIDs := []string{
        "bitcoin",
        "ethereum", 
        "tether",
        "binancecoin",
        "solana",
        "ripple",
        "cardano",
    }

	prices, err := client.GetPrices(coinIDs)
	if err != nil {
		log.Printf("❌ API error: %v", err)
		return
	}

	for coinName, data := range prices {
		currencyID, err := repo.GetCurrencyID(coinName)
		if err != nil {
			fmt.Printf("⚠️ Currency %s not found, skipping\n", coinName)
			continue
		}

		err = repo.SaveRate(models.ExchangeRate{
			CurrencyID: currencyID,
			Price:      data.USD,
		})

		if err != nil {
			fmt.Printf("❌ Failed to save %s: %v\n", coinName, err)
		} else {
			fmt.Printf("✅ %s: $%.2f\n", coinName, data.USD)
		}
	}

	fmt.Printf("✅ [%s] Rates updated\n", currentTime)
}
