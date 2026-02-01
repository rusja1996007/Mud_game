package main

import (
	"context"
	"fmt"
	"study/feature_postgres/simple_connect"
	"study/feature_postgres/simple_sql"
	"time"
)

func main() {

	ctx := context.Background()
	conn, err := simple_connect.CreateConnect(ctx)
	if err != nil {
		panic(err)
	}
	if err := simple_sql.CreateTable(ctx, conn); err != nil {
		panic(err)
	}
	tasks, err := simple_sql.SelectRows(ctx, conn)
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		if task.ID == 9 {
			task.Title = "Нюхать гея"
			task.Description = "Отсыпать насвая и оформить вкид"
			task.Completed = false
			n := time.Now()
			task.Completed_at = &n

			if err := simple_sql.UpdateTask(ctx, conn, task); err != nil {
				panic(err)
			}
			break

		}
	}

	fmt.Println("все кул дай новый  автомат е")

}
