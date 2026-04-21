package player

import (
	"Mud_game/Mud_Game/internal/domain/garden"
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
)

// личная зона каждого игрока
type PLayerZone struct {
	PlayerID   string
	HomeRoomID string
	GardenID   string //сад
	RoadID     string // дорога
	WellID     string
	Rooms      map[string]*room.Room
	Garden     *garden.Garden // ✅ добавляем структуру огорода
}

func CreatePlayerZone(playerID string, PlayerName string) *PLayerZone {

	//  Генерируем уникальные ID для комнат
	homeID := fmt.Sprintf("home_%s", playerID)
	gardenID := fmt.Sprintf("garden_%s", playerID)
	roadID := fmt.Sprintf("road_%s", playerID)
	wellID := fmt.Sprintf("well_%s", playerID)

	home := &room.Room{
		ID:          homeID,
		Name:        fmt.Sprintf("Дом %s", PlayerName),
		Description: "Твой уютный дом. Здесь можешь хранить вещи.",
		Exits: map[string]string{
			"road":   roadID,
			"garden": gardenID,
			"well":   wellID,
		},
		Items: []*item.ItemStack{
			item.GetItem("empty bottle", 3),
			item.GetItem("water bottle", 2),
			item.GetItem("tomato", 1),
			item.GetItem("potato", 1),
			item.GetItem("empty bag", 1),
			item.GetItem("leather bag", 1),
			item.GetItem("tomato seeds", 3),
			item.GetItem("potato seeds", 3),
			item.GetItem("burdock seeds", 1),
			item.GetItem("clover seeds", 1),
			item.GetItem("iron sword", 1),
			item.GetItem("knife", 1),
			item.GetItem("leather hood", 1),
			item.GetItem("leather armor", 1),
			item.GetItem("leather boots", 1),
			item.GetItem("wooden shield", 1),
			item.GetItem("silver ring", 1),
			item.GetItem("gold ring", 1),
			item.GetItem("black ring", 1),
			item.GetItem("burdock", 2),
			item.GetItem("clover", 2),
		},
	}

	road := &room.Room{
		ID:          roadID,
		Name:        "Тропинка от дома",
		Description: "Ты на тропинке,ведущая к твоему дому. Вдали виден город.",
		Exits: map[string]string{
			"garden": gardenID,
			"home":   homeID,
			"south":  "global_town", // пока так, потом заменим на настоящий ID города
		},
		Items: []*item.ItemStack{},
	}

	gardenRoom := &room.Room{
		ID:          gardenID,
		Name:        fmt.Sprintf("Огород %s", PlayerName),
		Description: "Твой маленький огород, можно выращивать растения",
		Exits: map[string]string{
			"road": roadID,
			"home": homeID,
		},
		Items: []*item.ItemStack{},
	}

	wellRoom := &room.Room{
		ID:          wellID,
		Name:        "Колодец",
		Description: "Старый каменный колодец с чистой водой",
		Exits: map[string]string{
			"home":   homeID,
			"garden": gardenID,
		},
		Items: []*item.ItemStack{},
	}

	// ✅ Создаем структуру Garden с 3 грядками
	playerGarden := garden.NewGarden(playerID, 3)

	// 3. Собираем все комнаты в карту
	rooms := map[string]*room.Room{
		homeID:   home,
		gardenID: gardenRoom,
		roadID:   road,
		wellID:   wellRoom,
	}
	// 4. Возвращаем зону
	return &PLayerZone{
		PlayerID:   playerID,
		HomeRoomID: homeID,
		GardenID:   gardenID,
		RoadID:     roadID,
		Rooms:      rooms,
		Garden:     playerGarden,
		WellID:     wellID,
	}

}
