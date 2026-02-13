DROP TABLE users;

/*migrate -path migrations -database "postgres://postgres:8626@localhost:5432/postgres?sslmode=disable" down 1-понизить на 1 версию, 
down-понизить на самую первую;
up 1 повысить на 1
up повысить на последнюю
version - проверить версию на данный момент
migrate create -ext sql -dir migrations -seq unique_phone  - создать миграцию с названием "unique_phone"
migrate -path migrations -database "postgres://postgres:8626@localhost:5432/postgres?sslmode=disable" force 3-принудительный скат на 3 версию
//sslmode=disable если используется везде и в го и в миграции прописываем