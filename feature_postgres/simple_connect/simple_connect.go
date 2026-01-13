package simple_connect

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// "postgres://UserName:password@HostName:5432/DataBaseName//
func CheckConnect() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:8626@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Подключение к БД прошло успешно")
}
