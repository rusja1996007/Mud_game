package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Открываем соединение
	//Стучится в дверь (gorm.Open)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к БД: %w", err)
	}
	//Если открыли - возвращает "пульт управления" (*gorm.DB)
	return db, nil
}
