package config

// Это как чертеж или форма. Мы говорим: "Конфигурация сервера состоит из этих трех вещей".
type ServerConfig struct {
	Port       int
	NameServer string
	MaxPlayers int
}

// Это как "заводские настройки". Значения, которые разумны по умолчанию.
func DefaultConfig() ServerConfig {
	return ServerConfig{
		Port:       4000,
		NameServer: "My_MUD_Server",
		MaxPlayers: 100,
	}
}
