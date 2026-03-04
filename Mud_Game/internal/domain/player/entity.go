package player

import "Mud_game/Mud_Game/internal/domain/item"

type Player struct {
	ID          string
	Name        string
	CurrentRoom string            //текущая комната
	Inventory   []*item.ItemStack // ← стопки предметов (название + кол-во)

}
