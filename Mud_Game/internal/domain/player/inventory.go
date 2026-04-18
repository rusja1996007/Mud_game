package player

import (
	"Mud_game/Mud_Game/internal/domain/item"
)

// /////////////////////////ДЛЯ ХЭНДЛЕРОВ///////////////////////////////////
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

// /////////////////////////////МЕТОДЫ СТРУКТУРЫ Player//////////////////////////////////
// возвращает количество занятых слотов
func (p *Player) GetUsedSlots() int {
	if p.Inventory == nil {
		return 0
	}
	return len(p.Inventory)
}

// максимальное количество слотов
func (p *Player) GetMaxSlots() int {
	baseSlots := 4 //базово 4 слота

	if p.Equipment == nil {
		return baseSlots
	}

	if p.Equipment.Bag != nil { //проверить, есть ли мешок
		baseSlots += p.Equipment.Bag.SlotBonus
	}
	return baseSlots
}

// возвращает свободны слоты
func (p *Player) GetFreeSlots() int {
	return p.GetMaxSlots() - p.GetUsedSlots()
}

// проверяет, есть ли место для нового предмета
func (p *Player) CanAddItem() bool {
	return p.GetFreeSlots() > 0
}

// добавить предмет в инвентарь
func (p *Player) AddItemToInventory(stack *item.ItemStack) bool {

	if !p.CanAddItem() {
		return false
	}

	if p.Inventory == nil {
		p.Inventory = []*item.ItemStack{}
	}

	for _, existing := range p.Inventory {
		if existing.Name == stack.Name {
			existing.Count += stack.Count
			return true
		}
	}
	newStack := &item.ItemStack{
		Name:          stack.Name,
		Count:         stack.Count,
		ItemType:      stack.ItemType,
		SlotBonus:     stack.SlotBonus,
		HungerRestore: stack.HungerRestore,
		ThirstRestore: stack.ThirstRestore,
	}
	p.Inventory = append(p.Inventory, newStack)

	return true
}

// находит индекс предмета по имени
// ВОЗВРАЩАЕТ:
// - индекс предмета (0, 1, 2...)
// - -1 если предмет не найден
func (p *Player) FindItemIndex(name string) int {
	for i, stack := range p.Inventory {
		if stack.Name == name {
			return i
		}

	}
	return -1
}
