package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	// Config хранит настройки подключения к БД
	//Это только то, что нужно для подключения к БД
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection создаёт подключение к PostgreSQL
func NewConnection(cfg Config) (*gorm.DB, error) {
	// Формируем строку подключения(Формирует адрес (dsn))

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Настройка логгера с увеличенным порогом SLOW SQL
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 1 * time.Second, // предупреждение только если запрос дольше 1 секунды
			LogLevel:      logger.Warn,     // логировать только ошибки и предупреждения
		},
	)

	// Открываем соединение
	//Стучится в дверь (gorm.Open)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к БД: %w", err)
	}
	//Если открыли - возвращает "пульт управления" (*gorm.DB)
	return db, nil
}
