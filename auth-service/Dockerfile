FROM golang:alpine AS builder

WORKDIR /app

COPY . .
# COPY auth-service ./auth-service
# COPY shared ./shared
# COPY go.work ./
RUN ls -la auth-service/

WORKDIR /app/auth-service

RUN go mod tidy && go build -o auth-service ./cmd/main.go

# RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-service/auth-service .
# COPY --from=builder /app/auth-service/internal/migrations ./migrations/
RUN mkdir -p /app/internal/migrations && chmod -R 755 /app/internal
CMD ["./auth-service"]  