package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// пришел хттп запрос
// взяли тело запроса
// преобразовали в структуру
// достали значение полей и передали в функцию InsertRow
// "SQL иньекция" - DROP DATABASE POSTGRES -ВСЯ БД Удалиться с данными, библиотека github.com/jackc/pgx/v5 защищает от этого
func InsertRow(ctx context.Context,
	conn *pgx.Conn,
	task TaskModel,
) error {
	sqlZapros := `
	INSERT INTO tasks(title, description, completed, created_at)
	VALUES ($1, $2, $3, $4);
	`
	_, err := conn.Exec(ctx,
		sqlZapros,
		task.Title,
		task.Description,
		task.Completed,
		task.Created_at)
	return err

}
