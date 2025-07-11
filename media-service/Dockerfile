FROM golang:alpine AS builder

WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

COPY . .

RUN ls -la media-service/

WORKDIR /app/media-service

RUN go mod tidy && go build -o media-service ./cmd/main.go

# RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/media-service/media-service .
# COPY --from=builder /app/auth-service/internal/migrations ./migrations/
CMD ["./media-service"]  