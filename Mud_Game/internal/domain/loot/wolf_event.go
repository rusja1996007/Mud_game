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
//*WolfFightResult

//string (сообщение о поломке)

//int (урон)

// int (опыт)
func FightWolf(weapon *item.ItemStack, tracking int, playerDefence int) (*WolfFightResult, string, int, int) {

	var winChance int

	if weapon == nil {
		winChance = 0
	} else {
		switch weapon.Name {
		case "knife":
			winChance = 20
		case "iron sword":
			winChance = 100 //////////////////////////////////////////временно
		default:
			winChance = 0
		}
	}

	//рол победы
	win := rand.Intn(100) < winChance

	//если проиграл:
	if !win {
		//волк наносит 5-15 урона
		wolfDamage := 5 + rand.Intn(11)

		return &WolfFightResult{
			Win:     false,
			Loot:    nil,
			Message: "Ты встретил волка! Ты попытался дать отпор, но понял, что не справляешься, и убежал.",
		}, "", wolfDamage, 0
	}

	//если вин генерируем дроп
	var loot []*item.ItemStack
	xp := 50
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

	//уменьшаем прочность оружия
	broken := false
	if weapon != nil {
		if weapon.Decrease(5) {
			broken = true
		}
	}

	result := &WolfFightResult{
		Win:     true,
		Loot:    loot,
		Message: "Ты встретил волка и одолел его!",
	}

	if broken {
		return result, "Твое оружие сломалось в бою!", 0, xp
	}

	return result, "", 0, xp

}
