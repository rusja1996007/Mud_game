package main

import (
	"Mud_game/Mud_Game/internal/delivery/tcp"
	"Mud_game/Mud_Game/internal/pkg/config"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"fmt"
	"strconv"
)

func main() {
	log := logger.NewSimpleLogger("MAIN") //log = видеокамеры наблюдения (записывают всё, что происходит)
	//Загружаем конфигурацию (порт, имя сервера и т.д.)
	cfg := config.DefaultConfig() //cfg = лицензия и документы (по каким правилам работаем)

	log.Info("Сервер запускается...")
	log.Info(fmt.Sprintf("Название: %s", cfg.NameServer))
	log.Info(fmt.Sprintf("Порт: %d", cfg.Port))
	//Создаем TCP сервер
	server := tcp.NewServer(strconv.Itoa(cfg.Port), log)
	go func() {
		//Запускаем сервер в отдельной горутине
		err := server.Start()
		if err != nil {
			log.Error("Ошибка:" + err.Error())
		}
	}()
	// 6. Блокируем main горутину навсегда
	// Если этого не сделать, программа завершится сразу после запуска
	// Пустой select блокируется навечно
	select {}

}
