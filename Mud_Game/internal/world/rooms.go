package world

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/room"
	"errors"
)

func InitRooms(repo room.Repository) error { //функция будет загружать все комнаты
	//Твой дом - первая комната
	home := &room.Room{
		ID:          "home_01",
		Name:        "Твой дома",
		Description: "Ты в своем скромном доме",
		Exits:       map[string]string{"south": "road_01"}, //road_01 куда может привести эта комната(твой дом) тоесть в дорогу
		Items: []item.ItemStack{
			{Name: "Empty bottle", Count: 20},
			{Name: "Empty bag", Count: 5},
		},
	}

	if err := repo.Save(home); err != nil { //Сохранить дом через repo.Save()
		return errors.New("Не получилось загрузить карту твоего дома" + err.Error())
	}

	//Дорога -вторая комната
	road := &room.Room{
		ID:          "road_01",
		Name:        "Дорога",
		Description: "Ты на дороге которая ведет  на юг",
		Exits:       map[string]string{"north": "home_01"}, //home_01 это куда может привести эта комната(дорога) тоесть домой
	}
	if err := repo.Save(road); err != nil { //Сохранить дорогу через repo.Save()
		return errors.New("Не получилось загрузить карту дороги")
	}
	return nil
}
