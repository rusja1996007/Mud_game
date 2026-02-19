package player

type Repository interface {
	// Save сохраняет игрока. Если игрок уже есть — обновляет.
	Save(player *Player) error
	FindByID(id string) (*Player, error)
	FindByName(name string) (*Player, error)
	Delete(id string) error
}

// КОНТРАКТ: если ты умеешь все выше - ты репозиторий.
