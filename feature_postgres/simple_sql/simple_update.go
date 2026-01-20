package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func UpdateTask(
	ctx context.Context,
	conn *pgx.Conn,
	task TaskModel,
) error {
	sqlZapros := `
	UPDATE tasks
	SET title=$1, description=$2, completed=$3, created_at=$4, completed_at=$5
	WHERE id=$6
	`
	_, err := conn.Exec(
		ctx,
		sqlZapros,
		task.Title,
		task.Description,
		task.Completed,
		task.Created_at,
		task.Completed_at,
		task.ID,
	)
	return err

}
