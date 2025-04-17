package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/JeanGrijp/learn-go-kafka/internal/handler"
	"github.com/JeanGrijp/learn-go-kafka/internal/kafka"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	slog.InfoContext(ctx, "main", "message", "Iniciando o servidor HTTP na porta 8080, http://localhost:8080/buscar-endereco?cep=12345678")

	brokers := []string{"kafka:9092"}
	slog.InfoContext(ctx, "main", "message", "Iniciando o Kafka com brokers", "brokers", brokers)

	err := kafka.InitProducer(brokers)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao iniciar o producer Kafka", "error", err)
		log.Fatalf("Erro ao iniciar o producer Kafka: %v", err)
	}

	slog.InfoContext(ctx, "main", "message", "Iniciando o async producer Kafka com brokers", "brokers", brokers)
	err = kafka.InitAsyncProducer(brokers)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao iniciar o async producer Kafka", "error", err)
		log.Fatalf("Erro ao iniciar o async producer Kafka: %v", err)
	}

	slog.InfoContext(ctx, "main", "message", "Iniciando o consumer Kafka com brokers", "brokers", brokers)
	producer, err := kafka.NewProducer()
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao criar producer Kafka", "error", err)
		log.Fatalf("Erro ao criar producer Kafka: %v", err)
	}
	slog.InfoContext(ctx, "main", "message", "Producer Kafka criado com sucesso")
	defer producer.Close()

	consumer, err := kafka.NewConsumer()
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao criar consumer Kafka", "error", err)
		log.Fatalf("Erro ao criar consumer Kafka: %v", err)
	}
	slog.InfoContext(ctx, "main", "message", "Consumer Kafka criado com sucesso")
	defer consumer.Close()

	slog.InfoContext(ctx, "main", "message", "Iniciando o consumer Kafka")
	go kafka.StartConsumer(ctx, "cep-topic")

	r := chi.NewRouter()
	r.Get("/buscar-endereco", handler.BuscarEnderecoHandler)

	log.Println("Servidor HTTP na porta 8080")
	slog.InfoContext(ctx, "main", "message", "Iniciando o servidor HTTP na porta 8080, http://localhost:8080/buscar-endereco?cep=12345678")
	http.ListenAndServe(":8080", r)

}
