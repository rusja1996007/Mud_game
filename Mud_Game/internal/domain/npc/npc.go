package npc

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type NPC struct {
	ID          string
	Name        string
	Description string
	RoomID      string
	Type        string         //("trader", "quest_giver", "trader_quest")//тип npc
	RefreshTime time.Duration  //время обновления
	LastRefresh time.Time      //последнее обновление
	Inventory   []*ItemForSale //товары если есть
	Quests      []string       // ID квеста если есть
	Mu          sync.RWMutex
}

type NPCModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	RoomID      string `gorm:"type:text"`
	Type        string `gorm:"type:text"`
	RefreshTime int    `gorm:"default:0"`
	LastRefresh int    `gorm:"default:0"`
	Inventory   string `gorm:"type:text"`
	Quests      string `gorm:"type:text"`
}

// структура предметов для продажи
type ItemForSale struct {
	ItemData   item.ItemData `json:"item_data"` //сам предмет
	Price      int           `json:"price"`     //цена
	Count      int           `json:"count"`     //кол-во
	BaseChance int           `json:"base_chance"`
}

// из БД в игру
func (m *NPCModel) ToEntity() (*NPC, error) {
	//парсим инвентарь
	var inventory []*ItemForSale

	if m.Inventory != "" {
		err := json.Unmarshal([]byte(m.Inventory), &inventory)
		if err != nil {
			return nil, errors.New("Не удалось преобразовать инвентарь npc в JSON")
		}
	}

	//парсим квесты
	var quests []string
	if m.Quests != "" {
		err := json.Unmarshal([]byte(m.Quests), &quests)
		if err != nil {
			return nil, errors.New("Не удалось преобразовать квесты")
		}
	}
	//Преобразовать RefreshTime
	rD := time.Duration(m.RefreshTime) * time.Second
	//Преобразовать LastRefresh
	lR := time.Unix(int64(m.LastRefresh), 0)
	//парсим торговца
	npc := &NPC{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		RoomID:      m.RoomID,
		Type:        m.Type,
		Inventory:   inventory,
		RefreshTime: rD,
		LastRefresh: lR,
		Quests:      quests,
	}
	return npc, nil
}

// из игры в БД
func FromEntity(n *NPC) (*NPCModel, error) {
	//преобразуем инвентарь
	inventJSON, err := json.Marshal(n.Inventory)
	if err != nil {
		return nil, errors.New("Не удалось преобразовать в JSON инвентарь")
	}

	//преобразуем квест
	questJSON, err := json.Marshal(n.Quests)
	if err != nil {
		return nil, errors.New("Не удалось преобразовать в JSON квесты")
	}
	//Преобразовать RefreshTime в int
	rTime := int(n.RefreshTime.Seconds())
	//Преобразовать LastRefresh в int
	lRefresh := n.LastRefresh.Unix()
	//возвращаем торговца
	traderModel := &NPCModel{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		RoomID:      n.RoomID,
		Type:        n.Type,
		Inventory:   string(inventJSON),
		RefreshTime: rTime,
		LastRefresh: int(lRefresh),
		Quests:      string(questJSON),
	}
	return traderModel, nil
}
