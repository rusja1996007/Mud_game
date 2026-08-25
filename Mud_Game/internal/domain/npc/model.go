package npc

import (
	"encoding/json"
	"errors"
	"time"
)

type TraderModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	RoomID      string `gorm:"type:text"`
	RefreshTime int    `gorm:"default:0"`
	LastRefresh int    `gorm:"default:0"`
	Inventory   string `gorm:"type:text"`
}

// из БД в игру
func (t *TraderModel) ToEntity() (*Trader, error) {
	//парсим инвентарь
	var inventory []*ItemForSale

	if t.Inventory != "" {
		err := json.Unmarshal([]byte(t.Inventory), &inventory)
		if err != nil {
			return nil, errors.New("Не удалось преобразовать инвентарь торговца в JSON")
		}
	}
	//Преобразовать RefreshTime
	rD := time.Duration(t.RefreshTime) * time.Second
	//Преобразовать LastRefresh
	lR := time.Unix(int64(t.LastRefresh), 0)
	//парсим торговца
	trader := &Trader{
		ID:              t.ID,
		Name:            t.Name,
		Description:     t.Description,
		RoomID:          t.RoomID,
		InventoryTrader: inventory,
		RefreshTime:     rD,
		LastRefresh:     lR,
	}
	return trader, nil
}

// из игры в БД
func FromEntity(t *Trader) (*TraderModel, error) {
	//преобразуем инвентарь
	inventJSON, err := json.Marshal(t.InventoryTrader)
	if err != nil {
		return nil, errors.New("Не удалось преобразовать в JSON инвентарь")
	}
	//Преобразовать RefreshTime в int
	rTime := int(t.RefreshTime.Seconds())
	//Преобразовать LastRefresh в int
	lRefresh := t.LastRefresh.Unix()
	//возвращаем торговца
	traderModel := &TraderModel{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		RoomID:      t.RoomID,
		Inventory:   string(inventJSON),
		RefreshTime: rTime,
		LastRefresh: int(lRefresh),
	}
	return traderModel, nil
}
