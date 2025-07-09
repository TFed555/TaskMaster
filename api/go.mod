module api

go 1.24.4

require auth-service v0.0.0

require (
	github.com/go-chi/chi/v5 v5.2.2 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	shared v0.0.0 // indirect
)

replace (
	auth-service => ../auth-service
	shared => ../shared
)
