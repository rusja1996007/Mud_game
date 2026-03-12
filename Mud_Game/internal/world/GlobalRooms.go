package world

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
)

// ✅ ДОБАВИТЬ после создания репозиториев
func InitGlobalTown(repo room.Repository) error {
	// Проверяем, есть ли уже город
	existing, _ := repo.FindByID("global_town")
	if existing != nil {
		fmt.Println("🌍 Город global_town загружен из БД")
		return nil // город уже есть
	}

	town := &room.Room{
		ID:          "global_town",
		Name:        "Городская площадь",
		Description: "Центральная площадь города. Сюда ведут дороги от домов игроков.",
		Exits:       map[string]string{}, // выходы будут добавляться позже
		Items:       []*item.ItemStack{},
	}
	if err := repo.Save(town); err != nil {
		return fmt.Errorf("Ошибка создания города: %v", err)
	}
	fmt.Println("🌍 Создан город global_town")
	return nil
}
