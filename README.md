# TaskMaster
Запуск отдельных сервисов:
cd ./media-service => go run ./cmd/main.go => тоже самое notes-service И auth-service
Запуск всего:
cd ./api => go run ./cmd/main.go

Пример запроса для получения тудушек
localhost:8082/api/todos?createdAt=2025-06-01T15:30:45Z&filter=before&offset=0&limit=1&title=Titl&category=d&filterStatus=completed
filterStatus = completed/current