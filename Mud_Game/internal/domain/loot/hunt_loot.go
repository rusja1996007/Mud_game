package loot

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"math/rand"
)

// LootItem представляет предмет с шансом выпадения
type LootItem struct {
	ItemData item.ItemData
	MinCount int //минимальное кол-во
	MaxCount int //макс. кол-во
	Chance   int //шанс в процентах(0-100)
}

// GenerateHuntLoot генерирует случайный лут на основе охоты
func GenerateHuntLoot(tracking int) []*item.ItemStack {
	var loot []*item.ItemStack

	for _, lootItem := range HuntLootTable {

		//базовый шанс
		chance := lootItem.Chance

		if tracking >= 5 {
			if lootItem.ItemData.Name == "coin" ||
				lootItem.ItemData.Name == "inonotus obliquus" ||
				lootItem.ItemData.Name == "cooper ring" ||
				lootItem.ItemData.Name == "rubroboletus satanas" {
				chance = chance * 130 / 100
			}
		}
		if rand.Intn(100) < chance {
			//выпало
			//расчет кол-ва
			count := lootItem.MinCount
			if lootItem.MaxCount > lootItem.MinCount {
				count += rand.Intn(lootItem.MaxCount - lootItem.MinCount + 1)
			}

			//добавление
			loot = append(loot, &item.ItemStack{
				Name:     lootItem.ItemData.Name,
				Count:    count,
				ItemType: lootItem.ItemData.ItemType,
			})
		}
	}
	return loot
}

// HuntLootTable - таблица лута для охоты
var HuntLootTable = []LootItem{
	{
		ItemData: item.ItemsDB["inonotus obliquus"], //чага
		MinCount: 1,
		MaxCount: 1,
		Chance:   10,
	},
	{
		ItemData: item.ItemsDB["clover"],
		MinCount: 1,
		MaxCount: 4,
		Chance:   50,
	},
	{
		ItemData: item.ItemsDB["broken sword"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   10,
	},
	{
		ItemData: item.ItemsDB["burdock"],
		MinCount: 1,
		MaxCount: 4,
		Chance:   50,
	},

	{
		ItemData: item.ItemsDB["hare meat"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   20,
	},

	{
		ItemData: item.ItemsDB["hare ears"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   10,
	},

	{
		ItemData: item.ItemsDB["hare paws"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   10,
	},

	{
		ItemData: item.ItemsDB["coin"],
		MinCount: 1,
		MaxCount: 2,
		Chance:   10,
	},

	{
		ItemData: item.ItemsDB["cooper ring"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   5,
	},

	{
		ItemData: item.ItemsDB["rubroboletus satanas"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   5,
	},

	{
		ItemData: item.ItemsDB["boletus edulis"],
		MinCount: 1,
		MaxCount: 1,
		Chance:   20,
	},

	{
		ItemData: item.ItemsDB["rubus caesius"], //ежевика
		MinCount: 1,
		MaxCount: 1,
		Chance:   50,
	},
}
