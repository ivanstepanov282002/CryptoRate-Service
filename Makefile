.PHONY: test test-cover test-unit test-integration bench build clean docker-build docker-push docker-up docker-down docker-test lint fmt

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

# Сборка Docker образов
docker-build:
	docker build -f Dockerfile.api -t cryptorate-api:latest .
	docker build -f Dockerfile.bot -t cryptorate-bot:latest .
	docker build -f Dockerfile.worker -t cryptorate-worker:latest .

# Загрузка Docker образов
docker-push:
	docker tag cryptorate-api:latest $(DOCKER_REGISTRY)/cryptorate-api:latest
	docker tag cryptorate-bot:latest $(DOCKER_REGISTRY)/cryptorate-bot:latest
	docker tag cryptorate-worker:latest $(DOCKER_REGISTRY)/cryptorate-worker:latest
	docker push $(DOCKER_REGISTRY)/cryptorate-api:latest
	docker push $(DOCKER_REGISTRY)/cryptorate-bot:latest
	docker push $(DOCKER_REGISTRY)/cryptorate-worker:latest

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

# Линтинг с автоматической установкой golangci-lint
lint:
	@which golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	golangci-lint run ./...

# Форматирование кода
fmt:
	go fmt ./...

# Запуск production compose
docker-up-prod:
	docker-compose -f docker-compose.prod.yml up -d

docker-down-prod:
	docker-compose -f docker-compose.prod.yml down

# Миграции БД
migrate-up:
	docker-compose exec postgres psql -U crypto_user -d crypto_db -f /docker-entrypoint-initdb.d/01-init.sql || echo "Migration file may not exist yet"

# Проверка безопасности
security-scan:
	@which gosec >/dev/null 2>&1 || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

# Генерация документации
docs:
	@which godoc >/dev/null 2>&1 || (echo "godoc not found, please install Go tools"; exit 1)
	godoc -http=:6060

# Установка всех инструментов разработки
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Быстрый линтинг для CI (без строгих проверок)
lint-ci:
	golangci-lint run --no-config --enable-only=govet,staticcheck ./...