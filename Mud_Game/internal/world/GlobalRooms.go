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

	//СОЗДАЕМ ГОРОД
	town := &room.Room{
		ID:          "global_town",
		Name:        "Городская площадь",
		Description: "Центральная площадь города. Сюда ведут дороги от домов игроков.",
		Exits: map[string]string{
			"hotel": "hotel",
		},
		Items: []*item.ItemStack{},
	}
	if err := repo.Save(town); err != nil {
		return fmt.Errorf("Ошибка создания города: %v", err)
	}
	fmt.Println("🌍 Создан город global_town")

	//СОЗДАЕМ ГОСТИНИЦУ
	hotel := &room.Room{
		ID:          "hotel",
		Name:        "Гостиница",
		Description: "Уютная гостиница с пропитанием. За 20 монет можно отдохнуть и восстановить все жизни, хорошо поесть и напиться. Время отдыха 5 минут. Используй команду <pay 20>",
		Exits: map[string]string{
			"town": "global_town", //выход из гостинницы в город
		},
		Items: []*item.ItemStack{},
	}
	if err := repo.Save(hotel); err != nil {
		return fmt.Errorf("Ошибка создания гостиницы: %v", err)
	}
	fmt.Println("🏨 Создана гостиница")

	//Добавляем выход из города в гостиницу
	town.Exits["hotel"] = "hotel"
	if err := repo.Save(town); err != nil {
		return fmt.Errorf("Ошибка обновления города: %v", err)
	}
	return nil
}
