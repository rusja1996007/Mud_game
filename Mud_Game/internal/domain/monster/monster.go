package monster

import (
	"time"
)

type Monster struct {
	ID            string
	Name          string
	Health        int
	MaxHealth     int
	MinDamage     int
	MaxDamage     int
	Defence       int
	FireDefence   int
	MagicDefence  int
	PoisonDefence int
	Experience    int
	RoomID        string //в какой комнате находится
	IsAlive       bool
	Description   string
	RespawnTime   time.Time
	TimeToLoot    time.Time //до какова времени можно обыскивать лут
	CastTime      int       //сколько ходов уже колдует(0-2)
	IsCasting     bool      //костует ли в этом ходу
	PoisonDamage  int       //урон от яда монстра
}

// создания гоблина в комнате
func NewGoblin(roomID string) *Monster {
	return &Monster{
		ID:            "goblin_1",
		Name:          "👹 Гоблин",
		Health:        100,
		MaxHealth:     100,
		MinDamage:     5,
		MaxDamage:     15,
		Defence:       5,
		FireDefence:   0,
		MagicDefence:  0,
		PoisonDefence: 0,
		Experience:    220,
		RoomID:        roomID,
		IsAlive:       true,
		Description:   "👹 Гоблин с кинжалом",
	}
}

func NewGoblinWarrior(roomID string) *Monster {
	return &Monster{
		ID:            "goblin_warrior",
		Name:          "👹 Гоблин мечник",
		Health:        130,
		MaxHealth:     130,
		MinDamage:     15,
		MaxDamage:     25,
		Defence:       10,
		FireDefence:   0,
		MagicDefence:  0,
		PoisonDefence: 5,
		Experience:    400,
		RoomID:        roomID,
		IsAlive:       true,
		Description:   "👹 Гоблин-воин с большим мечом",
	}
}

func NewGoblinShaman(roomID string) *Monster {
	return &Monster{
		ID:            "goblin_shaman",
		Name:          "🧙 Гоблин-шаман",
		Health:        20,
		MaxHealth:     20,
		MinDamage:     4,
		MaxDamage:     8,
		Defence:       0,
		FireDefence:   5,
		MagicDefence:  15,
		PoisonDefence: 10,
		Experience:    100,
		RoomID:        roomID,
		IsAlive:       true,
		Description:   "🧙 Гоблин-шаман с посохом",
		CastTime:      0,
		IsCasting:     true,
		PoisonDamage:  2,
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
