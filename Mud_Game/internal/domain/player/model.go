package player

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PlayerModel struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;size:50"`
	CurrentRoom string `gorm:"size:36"`
	Inventory   string `gorm:"type:text"`
}

//Визуализация таблицы в БД:
/*CREATE TABLE player_models (
    id SERIAL,
    name VARCHAR(50) UNIQUE,      -- короткое, с ограничением(уникальное, 50 символов)
    current_room VARCHAR(36),      -- короткое, с ограничением(ID комнаты)
    inventory TEXT,                 -- длинное, без ограничения( это JSON строка с предметами, много символов)
    created_at TIMESTAMP,
    updated_at TIMESTAMP,            -- это создаст GORM
    deleted_at TIMESTAMP
);
*/

// ToEntity конвертирует модель из БД в сущность Player
// *PlayerModel - это данные, которые пришли из базы данных
// *Player - возвращаем игрока из БД
// из БД в игру
func (m *PlayerModel) ToEntity() (*Player, error) {
	var inventory []*item.ItemStack

	//если инвентарь не пустой, парсим json
	if m.Inventory != "" {
		err := json.Unmarshal([]byte(m.Inventory), &inventory) //&inventory - это указатель на нашу коробку. Мы говорим: "положи результат вот в эту коробку(inventory)".
		if err != nil {
			return nil, errors.New("Не удалось преобразовать инвентарь из JSON")
		}
	}

	// Создаем и возвращаем игрока
	return &Player{
		ID:          fmt.Sprint(m.ID), //m.ID из БД превращаем в строку (fmt.Sprint)
		Name:        m.Name,           //m.Name из БД просто копируем
		CurrentRoom: m.CurrentRoom,    //m.CurrentRoom тоже копируем
		Inventory:   inventory,        //это наш преобразованный список предметов
	}, nil
}

// из игры в  БД:
func FromEntity(p *Player) (*PlayerModel, error) {
	// Превращаем стопки предметов в JSON строку для хранения в БД
	inventJSON, err := json.Marshal(p.Inventory)
	if err != nil {
		return nil, errors.New("Нe удалось преобразовать JSON в инвентарь ")
	}
	// Создаем модель БД из данных игрока
	return &PlayerModel{
		Name:        p.Name,
		CurrentRoom: p.CurrentRoom,
		Inventory:   string(inventJSON), // JSON нужно превратить в строку
		// ID: ?? пока не решаем
	}, nil
}
