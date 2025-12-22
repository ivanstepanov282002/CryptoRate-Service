.PHONY: test test-cover test-unit test-integration bench build clean

# Запуск всех тестов
test:
	go test ./... -v

# Запуск с покрытием
test-cover:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

# HTML отчет о покрытии
test-cover-html: test-cover
	go tool cover -html=coverage.out -o coverage.html

# Только unit тесты
test-unit:
	go test ./internal/... -v

# Бенчмарк тесты
bench:
	go test ./... -bench=. -benchmem

# Сборка проекта
build:
	go build ./cmd/...

# Очистка
clean:
	rm -f coverage.out coverage.html
	rm -f crypto-api crypto-bot crypto-worker

# Запуск в Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-test:
	docker-compose exec api go test ./...

# Линтинг
lint:
	golangci-lint run ./...

# Форматирование кода
fmt:
	go fmt ./...
