db_port := 8182

up:
	docker-compose up -d --build

down:
	docker-compose down

db:
	docker-compose up -d db

db-down:
	docker rm -f msg_db

migrate-new:
	goose -dir ./migrations create $(name) sql

migrate-up:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=dev host=localhost port=${db_port} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=dev host=localhost port=${db_port} sslmode=disable" down

send-msg:
	curl --header "Content-Type: application/json" \
           --request POST \
           --data '{"text":"Test message!"}' \
           http://localhost:8181/messages
