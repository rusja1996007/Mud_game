package player

import (
	"Mud_game/Mud_Game/internal/domain/garden"
	"Mud_game/Mud_Game/internal/domain/item"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PlayerModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"uniqueIndex;size:50"`
	CurrentRoom string `gorm:"size:36"` //где сейчас
	Inventory   string `gorm:"type:text"`
	Equipment   string `gorm:"type:text"`
	GardenData  string `gorm:"type:text"` // JSON с грядками("помидор: посажен в 14:00" )
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  gorm.DeletedAt `gorm:"index"`
}

// структура для JSON экипировки:
type EquipmentJSON struct {
	Weapon *item.ItemStack `json:"weapon"`
	Armor  *item.ItemStack `json:"armor"`
	Helmet *item.ItemStack `json:"helmet"`
	Bag    *item.ItemStack `json:"bag"`
	Shield *item.ItemStack `json:"shield"`
	Boots  *item.ItemStack `json:"boots"`
	Ring1  *item.ItemStack `json:"ring1"`
	Ring2  *item.ItemStack `json:"ring2"`
}

// Вспомогательные структуры для JSON
// GardenJSON - это инструкция, как записывать огород на бумажку:
type GardenJSON struct {
	PlayerID string      `json:"player_id"` //чей это огород
	Plots    []*PlotJSON `json:"plots"`     //список всех грядок
}

// 🌿 PlantJSON - как записывать растение:
type PlantJSON struct {
	Type      garden.PlantType `json:"type"`       //что это
	PlantedAt time.Time        `json:"planted_at"` //когда посадили
}

// PlotJSON - как записывать одну грядку:
type PlotJSON struct {
	ID    int        `json:"id"`              //грядка номер...
	Plant *PlantJSON `json:"plant,omitempty"` //omitempty - если поле пустое (nil для указателя), НЕ включай его в JSON//что здесь растет, если есть
}

// ToEntity конвертирует модель из БД в сущность Player
// *PlayerModel - это данные, которые пришли из базы данных
// *Player - возвращаем игрока из БД
// из БД в игру
func (m *PlayerModel) ToEntity() (*Player, error) {
	var inventory []*item.ItemStack

	//парсим инвентарь
	if m.Inventory != "" {
		err := json.Unmarshal([]byte(m.Inventory), &inventory) //&inventory - это указатель на нашу коробку. Мы говорим: "положи результат вот в эту коробку(inventory)".
		if err != nil {
			return nil, errors.New("Не удалось преобразовать инвентарь из JSON")
		}
	}

	//парсим огород
	var gardenObj *garden.Garden
	if m.GardenData != "" {
		var gardenJSON GardenJSON
		err := json.Unmarshal([]byte(m.GardenData), &gardenJSON)
		if err == nil { // Создаем новый огород
			gardenObj = garden.NewGarden(gardenJSON.PlayerID, len(gardenJSON.Plots))
			// Восстанавливаем растения на грядках
			for i, plotJSON := range gardenJSON.Plots {
				if plotJSON.Plant != nil {
					// Создаем растение с сохраненным временем посадки
					plant := garden.NewPlant(plotJSON.Plant.Type)
					plant.PlantedAt = plotJSON.Plant.PlantedAt
					gardenObj.Plots[i].Plant = plant
				}
			}
		}
	}

	equipment := &Equipment{}
	if m.Equipment != "" {
		var eqData EquipmentJSON
		if err := json.Unmarshal([]byte(m.Equipment), &eqData); err == nil {
			equipment = &Equipment{
				Weapon: eqData.Weapon,
				Armor:  eqData.Armor,
				Helmet: eqData.Helmet,
				Bag:    eqData.Bag,
				Shield: eqData.Shield,
				Boots:  eqData.Boots,
				Ring1:  eqData.Ring1,
				Ring2:  eqData.Ring2,
			}
		}
	}

	// Создаем и возвращаем игрока
	playerEntity := &Player{
		ID:          fmt.Sprint(m.ID), //m.ID из БД превращаем в строку (fmt.Sprint)
		Name:        m.Name,           //m.Name из БД просто копируем
		CurrentRoom: m.CurrentRoom,    //m.CurrentRoom тоже копируем
		Inventory:   inventory,        //это наш преобразованный список предметов
		Equipment:   equipment,
		Stats: &Stats{
			MaxSlots: 4,
			Hunger:   100,
			Thirst:   100,
			Level:    1,
		},
	}

	//// ✅ ВСЕГДА создаем Zone, даже если огород пустой!

	playerEntity.Zone = &PLayerZone{
		PlayerID: m.ID,
		Garden:   gardenObj,
		GardenID: fmt.Sprintf("garden_%s", m.ID), // генерируем ID
	}

	return playerEntity, nil
}

// из игры в  БД:
func FromEntity(p *Player) (*PlayerModel, error) {

	if p.ID == "" {
		return nil, errors.New("ID игрока пустой")
	}

	// Превращаем инвентарь в json
	inventJSON, err := json.Marshal(p.Inventory)
	if err != nil {
		return nil, errors.New("Нe удалось преобразовать JSON в инвентарь ")
	}

	// ✅ ОТЛАДКА: выводим JSON перед сохранением
	fmt.Printf("DEBUG Сохранение в БД: %s\n", string(inventJSON))

	//Превращаем огород в json
	var gardenJSON []byte
	if p.Zone != nil && p.Zone.Garden != nil {
		gardenData := GardenJSON{
			PlayerID: p.Zone.Garden.PlayerID,
			Plots:    make([]*PlotJSON, len(p.Zone.Garden.Plots)),
		}
		for i, plot := range p.Zone.Garden.Plots { //пробегаем по всем грядкам
			plotJSON := &PlotJSON{
				ID: plot.ID,
			}
			if plot.Plant != nil {
				plotJSON.Plant = &PlantJSON{
					Type:      plot.Plant.Type,
					PlantedAt: plot.Plant.PlantedAt,
				}
			}
			gardenData.Plots[i] = plotJSON
		}
		gardenJSON, _ = json.Marshal(gardenData)
	}

	var equipmentJSON []byte
	if p.Equipment != nil {
		eqData := EquipmentJSON{
			Weapon: p.Equipment.Weapon,
			Armor:  p.Equipment.Armor,
			Helmet: p.Equipment.Helmet,
			Bag:    p.Equipment.Bag,
			Shield: p.Equipment.Shield,
			Boots:  p.Equipment.Boots,
			Ring1:  p.Equipment.Ring1,
			Ring2:  p.Equipment.Ring2,
		}
		equipmentJSON, _ = json.Marshal(eqData)
	}

	// Создаем модель БД из данных игрока
	return &PlayerModel{
		ID:          p.ID,
		Name:        p.Name,
		CurrentRoom: p.CurrentRoom,
		Inventory:   string(inventJSON), // JSON нужно превратить в строку
		Equipment:   string(equipmentJSON),
		GardenData:  string(gardenJSON),
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}, nil
}
