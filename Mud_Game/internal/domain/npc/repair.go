package npc

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"time"
)

// npc держит отремонтированые вещи у себя чтобы ты забрал
type RepairOrder struct {
	PlayerID  string          `json:"player_id"`
	Item      *item.ItemStack `json:"item"`
	StartTime time.Time       `json:"start_time"` //когда начался ремонт
	EndTime   time.Time       `json:"end_time"`   //когда закончится ремонт
	Price     int             `json:"price"`      //цена ремонта
}

// расчет цены на ремонт
func CalculateRepairPrice(stack *item.ItemStack) int {

	if stack == nil || stack.Durability == 100 {
		return 0
	}

	HowMuchNeed := 100 - stack.Durability
	var rarityCoef float64 //коэф редкости

	switch stack.Rarity {
	case "", "common":
		rarityCoef = 0.6
	case "epic":
		rarityCoef = 1.7
	case "rare":
		rarityCoef = 1.2
	default:
		return 0
	}

	price := int(float64(HowMuchNeed) * 0.5 * rarityCoef)
	if price < 1 {
		price = 1
	}
	return price
}

// расчет времени на ремонт
func CalculateRepairTime(stack *item.ItemStack) time.Duration {

	if stack == nil || stack.Durability == 100 {
		return 0
	}

	HowMuchNeed := 100 - stack.Durability
	var rarityCoef float64 //коэф редкости

	switch stack.Rarity {
	case "", "common":
		rarityCoef = 0.6
	case "epic":
		rarityCoef = 1.3
	case "rare":
		rarityCoef = 0.9
	default:
		return 0
	}

	return time.Duration(float64(HowMuchNeed)*1.0*rarityCoef) * time.Second
}
