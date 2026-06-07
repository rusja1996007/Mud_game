package loot

import "Mud_game/Mud_Game/internal/domain/item"

// лут из пещеры
type CaveLootItem struct {
	ItemData      item.ItemData
	MinCount      int
	MaxCount      int
	BaseChance    int
	TrackingBonus int //бонус от следопытства
}

var CaveLootTable = []CaveLootItem{
	{
		ItemData:   item.ItemsDB["knife"],
		MinCount:   1,
		MaxCount:   1,
		BaseChance: 30,
	},
	{
		ItemData:   item.ItemsDB["coin"],
		MinCount:   1,
		MaxCount:   5,
		BaseChance: 100,
	},
	{
		ItemData:   item.ItemsDB["cooper ring"],
		MinCount:   1,
		MaxCount:   1,
		BaseChance: 5,
	},
}
