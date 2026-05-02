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
func GenerateHuntLoot(weapon *item.ItemStack, tracking int, defence int) ([]*item.ItemStack, *WolfFightResult, string, int) {
	var allLoot []*item.ItemStack
	var wolfResult *WolfFightResult
	var brokenMsg string
	var totalDamage int

	//1 - эвент с волком
	if rand.Intn(100) < 50 {

		result, msg, damage := FightWolf(weapon, tracking, defence)
		wolfResult = result
		if msg != "" {
			brokenMsg = msg
		}

		totalDamage = damage

		if wolfResult != nil && wolfResult.Win {
			allLoot = append(allLoot, wolfResult.Loot...)
		}
	}

	//2 - обычный лут
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
			allLoot = append(allLoot, &item.ItemStack{
				Name:     lootItem.ItemData.Name,
				Count:    count,
				ItemType: lootItem.ItemData.ItemType,
			})
		}
	}
	return allLoot, wolfResult, brokenMsg, totalDamage
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
