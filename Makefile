build:
	go build -o ./out/lowfoodmap-tg-bot cmd/main.go

run:
	go run ./cmd/main.go

start:
	docker-compose up -d

stop:
	docker-compose down

restart:
	make stop && make start