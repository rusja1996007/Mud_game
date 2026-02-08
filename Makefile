
include .env
export

service-run:
	@go run main.go
service-deploy:
	@docker compose up -d application
service-undeploy:
	@docker compose down application
migrate-up:
	migrate -path migrations -database ${CONN_STRING} up
migrate-down:
	migrate -path migrations -database ${CONN_STRING} down
migrate-sbros:
	migrate -path migrations -database ${CONN_STRING} force 1
migrate-version:
	migrate -path migrations -database ${CONN_STRING} version
