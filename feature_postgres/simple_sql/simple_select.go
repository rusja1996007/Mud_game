package simple_sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func SelectRows(ctx context.Context, conn *pgx.Conn) ([]TaskModel, error) {
	sqlZapros := `
	SELECT id,title,description, completed, created_at, completed_at
	FROM tasks
	ORDER BY id ASC
	`
	rows, err := conn.Query(ctx, sqlZapros)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]TaskModel, 0)
	for rows.Next() {
		var task TaskModel

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.Created_at,
			&task.Completed_at,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
		//printTaks(task)
	}
	return tasks, nil

}
func printTaks(
	task TaskModel) {
	fmt.Println("----------------------")
	fmt.Println("id :", task.ID)
	fmt.Println("title :", task.Title)
	fmt.Println("description :", task.Description)
	fmt.Println("completed :", task.Completed)
	fmt.Println("created_at :", task.Created_at)
	fmt.Println("completed_at :", task.Completed_at)
}
