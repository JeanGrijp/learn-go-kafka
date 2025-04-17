package handler

import (
	"log/slog"
	"net/http"

	"github.com/JeanGrijp/learn-go-kafka/internal/kafka"
)

func BuscarEnderecoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cep := r.URL.Query().Get("cep")

	slog.InfoContext(ctx, "BuscarEnderecoHandler", "message", "Iniciando o handler BuscarEnderecoHandler")

	slog.InfoContext(ctx, "BuscarEnderecoHandler", "cep", cep)

	if cep == "" {
		slog.ErrorContext(ctx, "BuscarEnderecoHandler", "message", "cep é obrigatório")
		http.Error(w, "cep é obrigatório", http.StatusBadRequest)
		return
	}

	// Enviar CEP para o Kafka de forma síncrona
	kafka.SendAsyncMessage(ctx, "cep-topic", cep)
	slog.InfoContext(ctx, "BuscarEnderecoHandler", "message", "CEP enviado para o Kafka (assíncrono)", "cep", cep)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("CEP enviado com sucesso! (assíncrono)"))
	slog.InfoContext(ctx, "BuscarEnderecoHandler", "message", "CEP enviado com sucesso! (assíncrono)", "cep", cep)
}
