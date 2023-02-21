build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down -v

run:
	go run main.go

enter:
	docker-compose exec maintenance bash
