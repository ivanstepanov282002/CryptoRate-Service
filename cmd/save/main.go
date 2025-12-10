package main

import (
	"CryptoRate-Service/internal/api"
	"CryptoRate-Service/internal/models"
	"CryptoRate-Service/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	//Подключение к БД
	connStr := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
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
	fmt.Println("✅ Connection to database")

	//Создание репозитория
	repo := repository.NewRepository(db)
	fmt.Printf("Repository created: %v\n", repo)

	//Получение цен из API
	client := api.NewCoinGeckoClient()
	prices, err := client.GetPrices([]string{"bitcoin", "ethereum"})
	if err != nil {
		log.Fatal("API error:", err)
	}

	//Сохранение цен в БД
	fmt.Println("\n💾 Сохранение в БД...")

	for coinName, data := range prices {
		currencyID, err := repo.GetCurrencyID(coinName)
		if err != nil {
			fmt.Printf("Currency %s not found in DB, skipping\n", coinName)
			continue
		}

		// Сохраняем курс
		err = repo.SaveRate(models.ExchangeRate{
			CurrencyID: currencyID,
			Price:      data.USD,
		})

		if err != nil {
			fmt.Printf("❌ Failed to save %s: %v\n", coinName, err)
		} else {
			fmt.Printf("✅ Saved %s: $%.2f\n", coinName, data.USD)
		}
	}

	fmt.Println("\n✅ Ready to save to database!")
}
