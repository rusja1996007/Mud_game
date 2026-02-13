package main

import (
	"Mud_game/Mud_Game/internal/pkg/config"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"fmt"
)

func main() {
	log := logger.NewSimpleLogger("MAIN")
	cfg := config.DefaultConfig()

	log.Info("Сервер запускается...")
	log.Info(fmt.Sprintf("Название: %s", cfg.NameServer))
	log.Info(fmt.Sprintf("Порт: %d", cfg.Port))

}
