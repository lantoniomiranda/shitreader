run:
	docker compose down -v
	docker compose up -d
	sleep 2
	go run cmd/main.go
