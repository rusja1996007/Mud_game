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
	GardenData  string `gorm:"type:text"` // JSON с грядками("помидор: посажен в 14:00" )
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  gorm.DeletedAt `gorm:"index"`
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

	// Создаем и возвращаем игрока
	playerEntity := &Player{
		ID:          fmt.Sprint(m.ID), //m.ID из БД превращаем в строку (fmt.Sprint)
		Name:        m.Name,           //m.Name из БД просто копируем
		CurrentRoom: m.CurrentRoom,    //m.CurrentRoom тоже копируем
		Inventory:   inventory,        //это наш преобразованный список предметов
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

	// Создаем модель БД из данных игрока
	return &PlayerModel{
		ID:          p.ID,
		Name:        p.Name,
		CurrentRoom: p.CurrentRoom,
		Inventory:   string(inventJSON), // JSON нужно превратить в строку
		GardenData:  string(gardenJSON),
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}, nil
}
