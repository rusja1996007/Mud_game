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
	Rooms      map[string]*room.Room
	Garden     *garden.Garden // ✅ добавляем структуру огорода
}

func CreatePlayerZone(playerID string, PlayerName string) *PLayerZone {

	//  Генерируем уникальные ID для комнат
	homeID := fmt.Sprintf("home_%s", playerID)
	gardenID := fmt.Sprintf("garden_%s", playerID)
	roadID := fmt.Sprintf("road_%s", playerID)

	home := &room.Room{
		ID:          homeID,
		Name:        fmt.Sprintf("Дом %s", PlayerName),
		Description: "Твой уютный дом. Здесь можешь хранить вещи.",
		Exits: map[string]string{
			"road":   roadID,
			"garden": gardenID,
		},
		Items: []*item.ItemStack{
			{Name: "empty bottle", Count: 20, ItemType: "container"},
			{Name: "empty bag", Count: 5, ItemType: "container"},
			{Name: "tomato seeds", Count: 5, ItemType: "seed"},
			{Name: "potato seeds", Count: 5, ItemType: "seed"},
			{Name: "meadow_clover seeds", Count: 5, ItemType: "seed"},
			{Name: "burdock seeds", Count: 5, ItemType: "seed"},
			{Name: "leather bag", Count: 1, ItemType: "bag", SlotBonus: 4}, //кож. сумка(для переноса увеличение инвентаря)
			{Name: "iron sword", Count: 1, ItemType: "weapon"},
			{Name: "leather hood", Count: 1, ItemType: "helmet"}, //кож. капюшон
			{Name: "leather armor", Count: 1, ItemType: "armor"}, //кож. броня
			{Name: "knife", Count: 1, ItemType: "weapon"},
			{Name: "silver ring", Count: 1, ItemType: "ring"},
			{Name: "gold ring", Count: 1, ItemType: "ring"},
			{Name: "black ring", Count: 1, ItemType: "ring"},
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

	// ✅ Создаем структуру Garden с 3 грядками
	playerGarden := garden.NewGarden(playerID, 3)

	// 3. Собираем все комнаты в карту
	rooms := map[string]*room.Room{
		homeID:   home,
		gardenID: gardenRoom,
		roadID:   road,
	}
	// 4. Возвращаем зону
	return &PLayerZone{
		PlayerID:   playerID,
		HomeRoomID: homeID,
		GardenID:   gardenID,
		RoadID:     roadID,
		Rooms:      rooms,
		Garden:     playerGarden,
	}

}
