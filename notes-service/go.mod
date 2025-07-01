module notes-service

go 1.24.4

require (
	github.com/go-chi/chi/v5 v5.2.2
	github.com/go-chi/cors v1.2.1
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/jmoiron/sqlx v1.4.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)

require shared v0.0.0

replace shared => ../shared
