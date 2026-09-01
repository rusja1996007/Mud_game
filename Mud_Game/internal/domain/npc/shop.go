package npc

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"math/rand"
	"time"
)

// СОздание торговца хламом
func NewJunkTrader() *NPC {
	return &NPC{
		ID:          "junk_trader",
		Name:        "jtr",
		Description: " Торговец путешественник, не знаешь что найдешь у него в продаже",
		RoomID:      "global_town",
		Inventory:   GenerateJunkItems(),
		RefreshTime: 30 * time.Second, ////////////для теста
		LastRefresh: time.Now(),
	}
}

// СОздание торговца оружием и броней
func NewWeaponTrader() *NPC {
	return &NPC{
		ID:          "weapon_trader",
		Name:        "wtr",
		Description: " Мастер своего дела, продающий качественное снаряжение",
		RoomID:      "global_town",
		Inventory:   GenerateWeaponItems(),
		RefreshTime: 30 * time.Second, /////////////тест
		LastRefresh: time.Now(),
	}
}

// //////////////////////////////////// товары торговца хламом//////////////////////////////
var JunkItems = []*ItemForSale{
	{
		ItemData:   item.ItemsDB["scroll blank"],
		Price:      10,
		Count:      3,
		BaseChance: 90,
	},
	{
		ItemData:   item.ItemsDB["black opal"],
		Price:      100,
		Count:      1,
		BaseChance: 20,
	},
	{
		ItemData:   item.ItemsDB["white opal"],
		Price:      20,
		Count:      1,
		BaseChance: 60,
	},
}

// генерация товаров для торговца хламом
func GenerateJunkItems() []*ItemForSale {
	result := []*ItemForSale{}

	for _, item := range JunkItems {
		if rand.Intn(100) < item.BaseChance {
			result = append(result, item)
		}
	}

	return result

}

/////////////////////////////////////товары торговца оружием и броней/////////////////////////////////////

var WeaponItems = []*ItemForSale{
	{
		ItemData:   item.ItemsDB["iron shield"],
		Price:      20,
		Count:      1,
		BaseChance: 70,
	},
	{
		ItemData:   item.ItemsDB["wooden spear"],
		Price:      18,
		Count:      1,
		BaseChance: 100,
	},
	{
		ItemData:   item.ItemsDB["knife"],
		Price:      5,
		Count:      1,
		BaseChance: 100,
	},
}

// генерация товаров для торговца оружием и броней
func GenerateWeaponItems() []*ItemForSale {
	result := []*ItemForSale{}

	for _, item := range WeaponItems {
		if rand.Intn(100) < item.BaseChance {
			result = append(result, item)
		}

	}
	return result
}
