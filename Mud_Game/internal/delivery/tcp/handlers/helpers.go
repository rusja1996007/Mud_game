package handlers

import "Mud_game/Mud_Game/internal/domain/item"

// addToInventory объединяет стопки в инвентаре
func AddToInventory(inventory []*item.ItemStack, newStack *item.ItemStack) []*item.ItemStack {
	if newStack == nil || newStack.Count == 0 {
		return inventory
	}
	for i, stack := range inventory {
		if stack.Name == newStack.Name {
			inventory[i].Count += newStack.Count
			return inventory
		}
	}
	return append(inventory, newStack)
}
