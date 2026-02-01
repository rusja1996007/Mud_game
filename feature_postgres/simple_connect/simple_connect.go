package simple_connect

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

// "postgres://UserName:password@HostName:5432/DataBaseName(postgres://postgres:8626@localhost:5432/postgres)//
func CreateConnect(ctx context.Context) (*pgx.Conn, error) {
	conStr := os.Getenv("CONN_STRING") //скрыли логин и пароль от БД(чувствительные данные), они на компе
	return pgx.Connect(ctx, conStr)

}
