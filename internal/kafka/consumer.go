package kafka

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/IBM/sarama"
	"github.com/JeanGrijp/learn-go-kafka/internal/model"
)

func StartConsumer(ctx context.Context, topic string) {
	slog.InfoContext(ctx, "Iniciando o consumidor Kafka", "topic", topic)
	consumer, err := sarama.NewConsumer(GetBrokers(), NewKafkaConfig())
	if err != nil {

		slog.ErrorContext(ctx, "Erro ao criar consumidor", "error", err)
		return
	}
	defer consumer.Close()
	slog.InfoContext(ctx, "Consumidor Kafka criado com sucesso", "topic", topic)

	part, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {

		slog.ErrorContext(ctx, "Erro ao consumir partição", "error", err)
		return
	}
	defer part.Close()
	slog.InfoContext(ctx, "Consumindo partição", "topic", topic, "partition", part.Messages())

	for {
		select {
		case msg := <-part.Messages():
			slog.InfoContext(ctx, "Mensagem recebida", "value", string(msg.Value))
			model.FetchViaCep(string(msg.Value))
			slog.InfoContext(ctx, "Mensagem processada", "value", string(msg.Value))
		case err := <-part.Errors():
			slog.ErrorContext(ctx, "Erro ao consumir", "error", err)
		case <-ctx.Done():
			slog.InfoContext(ctx, "Parando o consumidor Kafka", "topic", topic)
			return
		}
	}
}

func NewConsumer() (sarama.Consumer, error) {
	brokers := []string{getBrokerURL()}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Erro ao criar consumer: %v", err)
		return nil, err
	}

	return consumer, nil
}

func getBrokerURL() string {
	broker := os.Getenv("KAFKA_BROKER_URL")
	if broker == "" {
		broker = "kafka:9092"
	}
	return broker
}
