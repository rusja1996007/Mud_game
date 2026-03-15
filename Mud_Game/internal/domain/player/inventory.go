package player

import "Mud_game/Mud_Game/internal/domain/item"

// Has - проверяет, есть ли предмет в нужном количестве
func HasItem(inventory []*item.ItemStack, name string, count int) bool {
	for _, stack := range inventory {
		if stack.Name == name && stack.Count >= count {
			return true
		}
	}
	return false
}

// Remove - удаляет предметы из инвентаря
func RemoveItem(inventory *[]*item.ItemStack, name string, count int) bool {
	for i, stack := range *inventory {
		if stack.Name == name {
			if stack.Count > count {
				stack.Count -= count
				return true
			} else if stack.Count == count {
				// Убираем всю стопку
				*inventory = append((*inventory)[:i], (*inventory)[i+1:]...)
				return true
			}
		}
	}
	return false
}

// Add - добавляет предметы в инвентарь
func AddItem(inventory *[]*item.ItemStack, name string, count int) {
	for _, stack := range *inventory {
		if stack.Name == name {
			stack.Count += count
			return
		}
	}
	// Если нет такой стопки - создаем новую
	*inventory = append(*inventory, &item.ItemStack{
		Name:  name,
		Count: count,
	})

}
