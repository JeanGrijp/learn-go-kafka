.PHONY: kafka-up kafka-down kafka-restart docker-reload build run test lint

# Sobe os containers Kafka, Zookeeper e Kafka UI
kafka-up:
	docker-compose up -d --build

# Derruba os containers
kafka-down:
	docker-compose down

# Derruba e sobe novamente
kafka-restart: kafka-down kafka-up

docker-reload:
	docker-compose down
	docker-compose up -d --build
	docker-compose logs -f app

# Compila o binário da aplicação
build:
	go build -o bin/app ./cmd/main.go

# Roda a aplicação localmente
run:
	go run ./cmd/main.go

# Executa os testes
test:
	go test ./...

# Roda o linter (requer instalação do golangci-lint)
lint:
	golangci-lint run ./...