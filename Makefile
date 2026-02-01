
include .env
export

service-run:
	go run main.go
migrate-up:
	migrate -path migrations -database ${CONN_STRING} up
migrate-down:
	migrate -path migrations -database ${CONN_STRING} down
migrate-sbros:
	migrate -path migrations -database ${CONN_STRING} force 1
migrate-version:
	migrate -path migrations -database ${CONN_STRING} version
