package world

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/monster"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"time"
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
			"hotel":    "hotel",
			"dungeon":  "dungeon_entrance_goblins",
			"dungeon2": "dungeon_entrance_goblins_v2",
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

	//СОЗДАЕМ ВХОД В ПОДЗЕМЕЛЬЕ С ГОБЛИНОМ
	entrance := &room.Room{
		ID:          "dungeon_entrance_goblins",
		Name:        "Вход в подземелье",
		Description: "Тёмный каменный проход, ведущий в глубины.",
		Exits: map[string]string{
			"town": "global_town",    //выход на дорогу
			"down": "dungeon_goblin", //спуска дальше
		},
		Items:         []*item.ItemStack{},
		NextSpawnTime: time.Now(),
	}
	if err := repo.Save(entrance); err != nil {
		return fmt.Errorf("Ошибка создания входа в подземелье:%v", err)
	}

	//Создаем ПОДЗЕМЕЛЬЕ С гоблином
	goblinRoom := &room.Room{
		ID:          "dungeon_goblin",
		Name:        "Логово гоблинов",
		Description: "Маленькая пещера. ",
		Exits:       map[string]string{},
		Items:       []*item.ItemStack{},
		Monster:     monster.NewGoblin("dungeon_goblin"),
		ExitRoom:    "dungeon_entrance_goblins",
	}
	if err := repo.Save(goblinRoom); err != nil {
		return fmt.Errorf("Ошибка создания комнаты с гоблином:%v", err)
	}

	//СОЗДАЕМ ВХОД В ПОДЗЕМЕЛЬЕ С ДВУМЯ ГОБЛИНАМИ(ГЛУБОКАЯ ПЕЩЕРА)
	entrance2 := &room.Room{
		ID:          "dungeon_entrance_goblins_v2",
		Name:        "Вход в глубокую пещеру",
		Description: "Темный проход, ведущий в глубины.",
		Exits: map[string]string{
			"town": "global_town",
			"down": "dungeon_goblins_v2",
		},
		Items:         []*item.ItemStack{},
		NextSpawnTime: time.Now(),
	}
	if err := repo.Save(entrance2); err != nil {
		return fmt.Errorf("Ошибка создания комнаты входа в глубокую пещеру с гоблинами:%v", err)
	}

	//СОЗДАЕМ ПОДЗЕМЕЛЬЕ С ГОБЛинами
	goblinRoom2 := &room.Room{
		ID:          "dungeon_goblins_v2",
		Name:        "Логово двух гоблинов",
		Description: "Пещера, где два гоблина охраняют проход.",
		Exits:       map[string]string{},
		Items:       []*item.ItemStack{},
		MonsterS: []*monster.Monster{
			monster.NewGoblinWarrior("dungeon_goblins_v2"),
			monster.NewGoblinShaman("dungeon_goblins_v2"),
		},
		ExitRoom: "dungeon_entrance_goblins_v2",
	}
	if err := repo.Save(goblinRoom2); err != nil {
		return fmt.Errorf("Ошибка создания комнаты с двумя гоблинами:%v", err)
	}

	return nil

}
