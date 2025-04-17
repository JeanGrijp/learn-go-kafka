# Etapa 1: build da aplicação
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/main.go

# Etapa 2: imagem final enxuta
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
