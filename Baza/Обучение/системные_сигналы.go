package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("PID:", os.Getpid()) //получить ID запускаемого приложения

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	// ловим SIGTERM ИЛИ SIGINT и аккуратно заверошаем приложение
	fmt.Println("ДО отмены контекста...")
	<-ctx.Done()
	fmt.Println("ПОСЛЕ отмены контексты")

}
