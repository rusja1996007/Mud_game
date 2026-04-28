package loot

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"math/rand"
)

// WolfLootItem представляет предмет с волка
type WolfLootItem struct {
	Name     string
	ItemType string
	Chance   int
}

// весь возможный дроп с волка
var WolfLootTable = []WolfLootItem{

	{Name: "wolf meat", ItemType: "ingredients", Chance: 100},
	{Name: "wolf fang", ItemType: "ingredients", Chance: 20}, //клыки
	{Name: "wolf еуе", ItemType: "ingredients", Chance: 8},
	{Name: "wolf heart", ItemType: "ingredients", Chance: 5},
}

// результат боя с волком
type WolfFightResult struct {
	Win     bool
	Loot    []*item.ItemStack
	Message string
}

// симуляци боя
func FightWolf(weapon *item.ItemStack, tracking int) (*WolfFightResult, string) {

	var winChance int

	switch weapon.Name {
	case "knife":
		winChance = 20
	case "iron sword":
		winChance = 50
	default:
		winChance = 0
	}

	//рол победы
	win := rand.Intn(100) < winChance

	//если проиграл:
	if !win {
		return &WolfFightResult{
			Win:     false,
			Loot:    nil,
			Message: "Ты встретил волка! Ты попытался дать отпор, но понял, что не справляешься, и убежал.",
		}, ""
	}

	//если вин генерируем дроп
	var loot []*item.ItemStack

	for _, wolfItem := range WolfLootTable {
		chance := wolfItem.Chance

		//бонус от следопытства
		if tracking >= 8 {
			chance = chance * 150 / 100 //+50% шанс
		} else if tracking >= 5 {
			chance = chance * 120 / 100 //+20% шанс
		}

		if chance > 100 {
			chance = 100
		}

		if rand.Intn(100) < chance {
			loot = append(loot, &item.ItemStack{
				Name:     wolfItem.Name,
				Count:    1,
				ItemType: wolfItem.ItemType,
			})
		}
	}

	return &WolfFightResult{
		Win:     true,
		Loot:    loot,
		Message: "Ты встретил волка и одолел его!",
	}, ""
}
