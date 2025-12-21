package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cryptorate-service/internal/api/rest"
	"cryptorate-service/internal/repository"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к БД
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "crypto_user"),
		getEnv("POSTGRES_PASSWORD", "secure_password_123"),
		getEnv("POSTGRES_DB", "crypto_db"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}
	fmt.Println("✅ Connected to database")

	// Создаем репозиторий и хендлеры
	repo := repository.NewRepository(db)
	handler := rest.NewHandler(repo)

	// Настраиваем роутер
	router := mux.NewRouter()

	// Middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware) // Для веб-приложений

	// API Routes (версия 1)
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Курсы валют
	apiV1.HandleFunc("/rates", handler.GetRates).Methods("GET")
	apiV1.HandleFunc("/rates/{currency}", handler.GetRate).Methods("GET")
	apiV1.HandleFunc("/rates/{currency}/stats", handler.GetStats).Methods("GET")

	// Валюты
	apiV1.HandleFunc("/currencies", handler.GetCurrencies).Methods("GET")
	apiV1.HandleFunc("/currencies/{id}", handler.GetCurrencies).Methods("GET")

	// Системные
	apiV1.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Корневой маршрут
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
            "service": "Crypto Rates API",
            "version": "1.0.0",
            "endpoints": {
                "rates": "/api/v1/rates",
                "currency_stats": "/api/v1/rates/{currency}/stats",
                "currencies": "/api/v1/currencies",
                "health": "/api/v1/health"
            },
            "documentation": "/docs"
        }`)
	})

	// Настройка сервера
	port := getEnv("API_PORT", "8080")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second, //держит открытое соединение 60 секунд, если запросов нет, то закрывает его. При повторном запросе отсчет начинается с начала
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("🚀 API server started on http://localhost:%s\n", port)
		fmt.Printf("📚 API docs: http://localhost:%s/api/v1/rates\n", port)
		fmt.Printf("🏥 Health: http://localhost:%s/api/v1/health\n", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {//ИГнорируем ошибку если сервер был остановлен с помощью graceful shutdown
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала для остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	fmt.Println("👋 Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
