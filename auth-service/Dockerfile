FROM golang:alpine AS builder

WORKDIR /app


COPY . .
# COPY auth-service ./auth-service
# COPY shared ./shared
# COPY go.work ./
RUN ls -la auth-service/

# Переходим в директорию с main.go
WORKDIR /app/auth-service/cmd

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-service/auth-service .
# COPY --from=builder /app/auth-service/internal/migrations ./migrations/
RUN mkdir -p /app/internal/migrations && chmod -R 755 /app/internal
CMD ["./auth-service"]  