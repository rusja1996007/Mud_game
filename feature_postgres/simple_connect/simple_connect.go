package simple_connect

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// "postgres://UserName:password@HostName:5432/DataBaseName//
func CreateConnect(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, "postgres://postgres:8626@localhost:5432/postgres")

}
