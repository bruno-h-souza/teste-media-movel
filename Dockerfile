# Etapa 1: Construção
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copia e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário estaticamente
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api/main.go

# Etapa 2: Imagem de execução final
FROM alpine:latest  

WORKDIR /app
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]