FROM golang:alpine AS builder

WORKDIR /app
COPY . .

# Переходим в директорию с main.go
WORKDIR /app/auth-service/cmd

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /auth-service .
CMD ["./auth-service"]