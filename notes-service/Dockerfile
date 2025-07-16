FROM golang:alpine AS builder

WORKDIR /app

COPY . .


WORKDIR /app/notes-service

RUN go mod tidy && go build -o notes-service ./cmd/main.go

# RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/notes-service/notes-service .
# COPY --from=builder /app/auth-service/internal/migrations ./migrations/
CMD ["./notes-service"]  