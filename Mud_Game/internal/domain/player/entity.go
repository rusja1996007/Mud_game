package player

import "Mud_game/Mud_Game/internal/domain/item"

// // Equipment - экипировка (не занимает слоты рюкзака)
type Equipment struct {
	Weapon *item.ItemStack
	Armor  *item.ItemStack
	Helmet *item.ItemStack
	Bag    *item.ItemStack //мешок(сумка)
	Shield *item.ItemStack //щит
	Boots  *item.ItemStack //тапки
	Ring1  *item.ItemStack //кольцо1
	Ring2  *item.ItemStack //кольцо2

}

type Player struct {
	ID          string
	Name        string
	CurrentRoom string            //текущая комната
	Inventory   []*item.ItemStack // ← стопки предметов (название + кол-во)
	Equipment   *Equipment
	Zone        *PLayerZone
	Stats       *Stats // характеристики

}

type Stats struct {
	//БАЗОВЫЕ
	MaxSlots int //слоты все в рюкзаке
	Hunger   int //голод
	Thirst   int //жажда

	//производные(позже)
	Level      int
	Experience int
}
