package player

import (
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
			{Name: "Empty bottle", Count: 20},
			{Name: "Empty bag", Count: 5},
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

	garden := &room.Room{
		ID:          gardenID,
		Name:        fmt.Sprintf("Огород %s", PlayerName),
		Description: "Твой маленький огород, можно выращивать растения",
		Exits: map[string]string{
			"road": roadID,
			"home": homeID,
		},
		Items: []*item.ItemStack{
			{Name: "Griadka", Count: 3},
		},
	}
	// 3. Собираем все комнаты в карту
	rooms := map[string]*room.Room{
		homeID:   home,
		gardenID: garden,
		roadID:   road,
	}
	// 4. Возвращаем зону
	return &PLayerZone{
		PlayerID:   playerID,
		HomeRoomID: homeID,
		GardenID:   gardenID,
		RoadID:     roadID,
		Rooms:      rooms,
	}

}
