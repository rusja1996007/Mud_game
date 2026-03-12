package main

import (
	"Mud_game/Mud_Game/internal/delivery/tcp"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
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

	"gorm.io/gorm"
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

	// 3️⃣ ОБЪЯВЛЯЕМ ПЕРЕМЕННые ДЛЯ РЕПОЗИТОРИЯ нужна для запуска сервера

	var pRepo player.Repository
	var rRepo room.Repository

	var database *gorm.DB
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
		database, err = db.NewConnection(dbConfig)
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
	// ===== ВЫБОР РЕПОЗИТОРИЯ ДЛЯ КОМНАТ =====
	if cfg.UsePostgres {
		// Используем PostgreSQL для комнат
		rRepo, err = roomRepo.NewPostgresRepository(database)
		if err != nil {
			log.Error("Ошибка создания репозитория комнат: " + err.Error())
			return
		}
		log.Info("✅ Комнаты будут сохраняться в PostgreSQL")
	} else {
		//// Используем in-memory
		rRepo = roomRepo.NewMemoryRepository()
		log.Info("📝 Комнаты в памяти (не сохраняются)")
	}

	//// ===== ИНИЦИАЛИЗАЦИЯ МИРА(общего города) =====
	if err := world.InitGlobalTown(rRepo); err != nil {
		log.Error("Ошибка создания города: " + err.Error())
		return
	}

	log.Info("Сервер запускается...")
	log.Info(fmt.Sprintf("Название: %s", cfg.NameServer))
	log.Info(fmt.Sprintf("Порт: %d", cfg.Port))

	//создание сервера
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
