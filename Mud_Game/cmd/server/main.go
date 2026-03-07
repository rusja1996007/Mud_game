package main

import (
	"Mud_game/Mud_Game/internal/delivery/tcp"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/pkg/config"
	"Mud_game/Mud_Game/internal/pkg/db"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"Mud_game/Mud_Game/internal/repository/memoryrepo"
	playerRepo "Mud_game/Mud_Game/internal/repository/player" // ← алиас! чтобы не было конфликта имён
	roomRepo "Mud_game/Mud_Game/internal/repository/room"     // ← алиас!
	"Mud_game/Mud_Game/internal/world"
	"flag"
	"fmt"
	"strconv"
)

func main() {
	log := logger.NewSimpleLogger("MAIN") //log = видеокамеры наблюдения (записывают всё, что происходит)
	//Загружать конфиг из файла
	configPath := flag.String("config", "config.yaml", "Путь к конфигу")
	flag.Parse() // парсим аргументы командной строки

	// Загружаем конфиг
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Error("Ошибка загрузки конфига:" + err.Error())
		return
	}

	// 3️⃣ ОБЪЯВЛЯЕМ ПЕРЕМЕННУЮ ДЛЯ РЕПОЗИТОРИЯ нужна для запуска сервера
	// player.Repository - это ИНТЕРФЕЙС (может быть хоть memory, хоть postgres)
	var pRepo player.Repository

	// Создаём конфиг для БД из данных, которые загрузили
	// ВЫБИРАЕМ, КАКОЙ РЕПОЗИТОРИЙ ИСПОЛЬЗОВАТЬ
	if cfg.UsePostgres {
		// ========== ВЕТКА POSTGRESQL ==========
		dbConfig := db.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			DBName:   cfg.Database.DBName,
			SSLMode:  cfg.Database.SSLMode,
		}

		//Подключаемся к БД
		database, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Error("Не удалось подключиться к БД:" + err.Error())
			return
		}
		log.Info("✅ Подключено к PostgreSQL")

		//	Создаём PostgreSQL репозиторий
		pRepo, err = playerRepo.NewPostgresRepository(database)
		if err != nil {
			log.Error("Ошибка создания репозитория: " + err.Error())
			return
		}
	} else {
		// ========== ВЕТКА IN-MEMORY ==========
		pRepo = memoryrepo.NewMemoryRepository()
		log.Info("📝 Используется in-memory хранилище")
	}

	log.Info("Сервер запускается...")
	log.Info(fmt.Sprintf("Название: %s", cfg.NameServer))
	log.Info(fmt.Sprintf("Порт: %d", cfg.Port))

	//Создаем репозиторий комнат (rRepo)+ создание сервера
	rRepo := roomRepo.NewMemoryRepository()
	//Загружаем комнаты
	if err := world.InitRooms(rRepo); err != nil {
		log.Error("Ошибка загрузки комнат :" + err.Error())
		return
	}

	server := tcp.NewServer(strconv.Itoa(cfg.Port), log, pRepo, rRepo)
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

//команда для подклюючения к серверу telnet localhost 4000
