package monster

import (
	"time"
)

type Monster struct {
	ID          string
	Name        string
	Health      int
	MaxHealth   int
	MinDamage   int
	MaxDamage   int
	Defence     int
	Experience  int
	RoomID      string //в какой комнате находится
	IsAlive     bool
	Description string
	RespawnTime time.Time
	TimeToLoot  time.Time //до какова времени можно обыскивать лут
}

// создания гоблина в комнате
func NewGoblin(roomID string) *Monster {
	return &Monster{
		ID:          "goblin_1",
		Name:        "Гоблин",
		Health:      100,
		MaxHealth:   100,
		MinDamage:   5,
		MaxDamage:   15,
		Defence:     5,
		Experience:  220,
		RoomID:      roomID,
		IsAlive:     true,
		Description: "👹 Гоблин с кинжалом",
	}
}

// проверяем надо ли возрождать монстра
func (m *Monster) CheckRespawn() bool {
	if !m.IsAlive && time.Now().After(m.RespawnTime) {
		m.IsAlive = true
		m.Health = m.MaxHealth
		m.RespawnTime = time.Time{}
		return true // респавн произошёл
	}
	return false
}
