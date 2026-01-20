package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DeleteRow(ctx context.Context, conn *pgx.Conn, tasksID []int) error {
	sqlZapros := `
	DELETE FROM tasks
	WHERE id = ANY($1);
	`
	_, err := conn.Exec(ctx, sqlZapros, tasksID)
	return err

}
