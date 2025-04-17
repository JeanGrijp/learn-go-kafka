package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
)

func NewProducer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(GetBrokers(), NewKafkaConfig())
}

func PublishMessage(producer sarama.SyncProducer, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := producer.SendMessage(msg)
	return err
}

var Producer sarama.SyncProducer
var AsyncProducer sarama.AsyncProducer

func InitProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	prod, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	Producer = prod
	return nil
}

func InitAsyncProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	prod, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return err
	}

	AsyncProducer = prod
	go handleAsyncProducerResponses()
	return nil
}

func handleAsyncProducerResponses() {
	for {
		select {
		case success := <-AsyncProducer.Successes():
			slog.Info("Mensagem enviada para o Kafka (assíncrono)",
				"topic", success.Topic,
				"partition", success.Partition,
				"offset", success.Offset,
			)
		case err := <-AsyncProducer.Errors():
			slog.Error("Erro ao enviar mensagem Kafka (assíncrono)", "err", err)
		}
	}
}

func SendMessage(ctx context.Context, topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := Producer.SendMessage(msg)
	if err != nil {
		slog.ErrorContext(ctx, "Erro ao enviar mensagem Kafka", "err", err)
		return err
	}

	slog.InfoContext(ctx, "Mensagem enviada para o Kafka",
		"topic", topic,
		"partition", partition,
		"offset", offset,
	)
	return nil
}

func SendAsyncMessage(ctx context.Context, topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	AsyncProducer.Input() <- msg
}
