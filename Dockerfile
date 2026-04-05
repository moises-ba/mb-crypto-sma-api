# Estágio de Compilação
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro (otimiza o cache das camadas)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário a partir do caminho que você especificou
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Estágio Final (Imagem de execução)
FROM alpine:latest

WORKDIR /root/

# Copia apenas o binário do estágio anterior
COPY --from=builder /app/main .

# Exponha a porta que o Gin está usando
EXPOSE 8080

CMD ["./main"]