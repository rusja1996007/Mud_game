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

	return len(p.Inventory) + len(p.Equipment.BagItems)
}

// максимальное количество слотов
func (p *Player) GetMaxSlots() int {
	baseSlots := 8 //базово 8 слота
	if p.Equipment.Bag != nil {
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

func (p *Player) AddItemToInventory(stack *item.ItemStack) bool {
	// 1. Сначала пытаемся добавить в инвентарь
	if len(p.Inventory) < 8 {
		for _, existing := range p.Inventory {
			if existing.CanStackWith(stack) {
				existing.Count += stack.Count
				return true
			}
		}
		p.Inventory = append(p.Inventory, stack)
		return true
	}

	// 2. Если инвентарь полон — пробуем в мешок
	if p.Equipment.Bag != nil {
		bagMax := p.Equipment.Bag.SlotBonus
		if len(p.Equipment.BagItems) < bagMax {
			for _, existing := range p.Equipment.BagItems {
				if existing.CanStackWith(stack) {
					existing.Count += stack.Count
					return true
				}
			}
			p.Equipment.BagItems = append(p.Equipment.BagItems, stack)
			return true
		}
	}

	return false // нет места
}

// находит индекс предмета по имени
// ВОЗВРАЩАЕТ:
// - индекс предмета (0, 1, 2...)
// - -1 если предмет не найден
// FindItemIndex находит индекс предмета по имени (только в инвентаре)
func (p *Player) FindItemIndex(name string) int {
	for i, stack := range p.Inventory {
		if stack.Name == name {
			return i
		}
	}
	return -1
}

// FindItemInBag находит индекс предмета по имени в мешке
func (p *Player) FindItemInBag(name string) int {
	for i, stack := range p.Equipment.BagItems {
		if stack.Name == name {
			return i
		}
	}
	return -1
}

// возвращает предмет по номеру и его индекс
func (p *Player) FindItemByNumber(number int) (*item.ItemStack, int) {
	if number < 1 {
		return nil, -1
	}

	// Инвентарь (1-8)
	if number <= len(p.Inventory) {
		return p.Inventory[number-1], number - 1
	}

	// Мешок (9-12)
	bagStart := len(p.Inventory)
	bagIndex := number - bagStart - 1
	if bagIndex >= 0 && bagIndex < len(p.Equipment.BagItems) {
		return p.Equipment.BagItems[bagIndex], bagIndex
	}

	return nil, -1
}

// FindItemGlobalByName ищет предмет по имени во всех хранилищах
func (p *Player) FindItemGlobalByName(name string) (int, bool) {
	// Инвентарь
	for i, stack := range p.Inventory {
		if stack.Name == name {
			return i, false
		}
	}
	// Мешок
	for i, stack := range p.Equipment.BagItems {
		if stack.Name == name {
			return i, true
		}
	}
	return -1, false
}

// RemoveItemFromStorage удаляет предмет из инвентаря или мешка
func (p *Player) RemoveItemFromStorage(name string, inBag bool, index int) {
	if inBag {
		p.Equipment.BagItems = append(p.Equipment.BagItems[:index], p.Equipment.BagItems[index+1:]...)
	} else {
		p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	}
}

// RemoveOneItem удаляет 1 единицу предмета из инвентаря или мешка
func (p *Player) RemoveOneItem(name string, inBag bool, index int) {
	var stack *item.ItemStack
	if inBag {
		stack = p.Equipment.BagItems[index]
	} else {
		stack = p.Inventory[index]
	}

	if stack.Count > 1 {
		stack.Count--
	} else {
		// Если была 1 штука — удаляем стопку
		if inBag {
			p.Equipment.BagItems = append(p.Equipment.BagItems[:index], p.Equipment.BagItems[index+1:]...)
		} else {
			p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
		}
	}
}
