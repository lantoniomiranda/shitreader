run:
	docker compose down -v
	docker compose up -d
	go run cmd/main.go
