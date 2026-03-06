package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Это как чертеж или форма. Мы говорим: "Конфигурация сервера состоит из этих  вещей".
type ServerConfig struct {
	Port       int    `yaml:"port"`
	NameServer string `yaml:"name_server"`
	MaxPlayers int    `yaml:"max_players"`
	// Добавляем настройки БД
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`  //имя базы данных
		SSLMode  string `yaml:"sslmode"` //режим шифрования подключения к БД 🔐
	} `yaml:"database"`
	UsePostgres bool `yaml:"use_postgres"` // переключатель
}

// Это как "заводские настройки". Значения, которые разумны по умолчанию.
func DefaultConfig() ServerConfig {
	cfg := ServerConfig{
		Port:        4000,
		NameServer:  "My_MUD_Server",
		MaxPlayers:  100,
		UsePostgres: false, //пока будем использовать память по умолчанию
	}
	// Значения по умолчанию для БД
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "muduser"
	cfg.Database.Password = "mudpassword"
	cfg.Database.DBName = "mudgame"
	cfg.Database.SSLMode = "disable" //disable - без шифрования,require - требовать шифрование,verify-full - максимальная безопасность

	return cfg
}

//# В psql мы выполнили:
//CREATE DATABASE mudgame;  // ← создали базу с именем "mudgame"
//В конфиге указываем:
//database:
//dbname: "mudgame"   ← говорим игре: "подключайся к БД с именем mudgame"

func Load(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path) //Берёт путь к файлу и читает всё содержимое
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг %s: %w", path, err)
	}
	var cfg ServerConfig //Создаём пустую структуру ServerConfig, куда положим распарсенные данные
	//Берёт сырые данные из файла (data)
	//Превращает YAML в Go-структуру
	//Кладёт результат в &cfg (в нашу пустую структуру)
	if err := yaml.Unmarshal(data, &cfg); err != nil { //Парсим YAML в структуру
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}
	return &cfg, nil

}
