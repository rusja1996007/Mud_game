package loot

import "Mud_game/Mud_Game/internal/domain/item"

// лут из ....
type CaveLootItem struct {
	ItemData      item.ItemData
	MinCount      int
	MaxCount      int
	BaseChance    int
	TrackingBonus int //бонус от следопытства
}

// лут из пещеры c 1 гоблином
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

// лут из пущеры с 2 гоблинами
var CaveV2LootTable = []CaveLootItem{
	{
		ItemData:   item.ItemsDB["coin"],
		MinCount:   5,
		MaxCount:   20,
		BaseChance: 20,
	},
	{
		ItemData:   item.ItemsDB["scroll fireball"],
		MinCount:   1,
		MaxCount:   1,
		BaseChance: 100,
	},
}

// лут из глубин(после пешеры с 2 гоблинами)
var GLubiniLootTable = []CaveLootItem{
	{
		ItemData:   item.ItemsDB["coin"],
		MinCount:   1,
		MaxCount:   5,
		BaseChance: 40,
	},
	{
		ItemData:   item.ItemsDB["scroll heal"],
		MinCount:   1,
		MaxCount:   1,
		BaseChance: 100,
	},
}
