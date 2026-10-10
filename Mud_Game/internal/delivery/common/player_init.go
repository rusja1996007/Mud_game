package common

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"time"
)

//файлы, общие для TCP и WebSocket.

func InitPlayer(name string, playerRepo player.Repository, roomRepo room.Repository) (*player.Player, error, string) {
	//Найти игрока в БД
	currentPlayer, err := playerRepo.FindByName(name)
	if err != nil {
		return nil, err, ""
	}

	if currentPlayer != nil {
		return currentPlayer, nil, fmt.Sprintf("С возвращением, %s\n", name)
	}

	//аналогично как в TCP
	// 1. Генерируем ID
	id := fmt.Sprintf("player_%d", time.Now().UnixNano())

	// 2. Создаём зону
	zone := player.CreatePlayerZone(id, name)

	// 3. Сохраняем все комнаты зоны
	for _, r := range zone.Rooms {
		if err := roomRepo.Save(r); err != nil {
			return nil, err, ""
		}
	}

	// 4. Связываем с городом (global_town)
	townInterface, err := roomRepo.FindByID("global_town")
	if err != nil {
		return nil, err, ""
	}

	townRoom, ok := townInterface.(*room.Room)
	if !ok {
		return nil, err, ""
	}

	nameExit := fmt.Sprintf("дом %s", name)
	townRoom.Exits[nameExit] = zone.RoadID
	townRoom.TownExits = append(townRoom.TownExits, room.TownExit{
		Name:    nameExit,
		RoomID:  zone.RoadID,
		OwnerID: id,
	})

	if err := roomRepo.Save(townRoom); err != nil {
		return nil, err, ""
	}

	// 5. Создаём игрока
	strength := 3
	newPlayer := &player.Player{
		ID:          id,
		Name:        name,
		CurrentRoom: zone.HomeRoomID,
		Inventory:   []*item.ItemStack{},
		Equipment:   &player.Equipment{},
		Stats: &player.Stats{
			MaxSlots:   8,
			Hunger:     100,
			Thirst:     100,
			Health:     50 + 5*strength,
			Strength:   3,
			Dexterity:  3,
			Intelect:   2,
			Tracking:   6,
			Level:      1,
			Experience: 0,
		},
		Zone: zone,
	}

	// 6. Сохраняем
	if err := playerRepo.Save(newPlayer); err != nil {
		return nil, err, ""
	}

	return newPlayer, nil, fmt.Sprintf("Привет %s! Добро пожаловать в игру!\n", name)

}
