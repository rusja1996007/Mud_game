package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// надеть...
func HandleWear(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//     Убираем "wear " из команды

	args, found := strings.CutPrefix(cmd, "wear ")
	if !found {
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		fmt.Fprintf(conn, "Что надеть? Использование: wear <предмет>\n> ")
		return
	}

	//название предмета(второй аргумент)
	itemName := args // полное название предмета (например "iron sword")

	//Найди предмет в инвентаре
	var index int = -1
	var inBag bool

	if num, err := strconv.Atoi(itemName); err == nil {
		target, idx := p.FindItemByNumber(num)
		if target == nil {
			fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", num)
			return
		}
		index = idx
		itemName = target.Name
		//где лежит предмет:
		if num <= len(p.Inventory) {
			inBag = false
		} else {
			inBag = true
		}
	} else {
		index, inBag = p.FindItemGlobalByName(itemName)
	}

	if index == -1 {
		fmt.Fprintf(conn, "У тебя нет предмета %s в инвентаре\n> ", itemName)
		return
	}
	//Получи предмет из инвентаря
	var it *item.ItemStack
	if inBag {
		it = p.Equipment.BagItems[index]
	} else {
		it = p.Inventory[index]
	}

	//Еда, напитки, семена, материалы,свитки
	if it.ItemType == "food" || it.ItemType == "drink" || it.ItemType == "seed" || it.ItemType == "material" || it.ItemType == "container" || it.ItemType == "scroll" {
		fmt.Fprintf(conn, "Этот предмет нельзя надеть.\n> ")
		return
	}

	// ✅ Проверка: не сломан ли предмет
	if it.Durability <= 0 {
		fmt.Fprintf(conn, "Этот предмет сломан, его нельзя надеть.\n> ")
		return
	}

	switch it.ItemType {

	//оружие
	case "weapon":
		if p.Equipment.Weapon != nil {
			p.AddItemToInventory(p.Equipment.Weapon)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Weapon))
		}
		p.Equipment.Weapon = it // ← добавить!
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот оружия\n> ", item.GetColoredName(it))
	//броня
	case "armor":
		if p.Equipment.Armor != nil {
			p.AddItemToInventory(p.Equipment.Armor)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Armor))
		}
		//вооружаем
		p.Equipment.Armor = it
		//удаляем из инвентаря
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот брони\n> ", item.GetColoredName(it))
	//шлем
	case "helmet":
		if p.Equipment.Helmet != nil {
			p.AddItemToInventory(p.Equipment.Helmet)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Helmet))
		}
		//вооружаем
		p.Equipment.Helmet = it
		//удаляем из инвентаря
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот шлема\n> ", item.GetColoredName(it))
	//ботинки
	case "boots":
		if p.Equipment.Boots != nil {
			p.AddItemToInventory(p.Equipment.Boots)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Boots))
		}
		//вооружаем
		p.Equipment.Boots = it
		//удаляем из инвентаря
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот обуви\n> ", item.GetColoredName(it))
	//щит
	case "shield":
		if p.Equipment.Shield != nil {
			p.AddItemToInventory(p.Equipment.Shield)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Shield))
		}
		//вооружаем
		p.Equipment.Shield = it
		//удаляем из инвентаря
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот щита\n> ", item.GetColoredName(it))
	//сумка
	case "bag":
		if p.Equipment.Bag != nil {
			p.AddItemToInventory(p.Equipment.Bag)
			fmt.Fprintf(conn, "Ты снял %s\n", item.GetColoredName(p.Equipment.Bag))
		}
		//вооружаем
		p.Equipment.Bag = it
		//удаляем из инвентаря
		p.RemoveItemFromStorage(itemName, inBag, index)
		fmt.Fprintf(conn, "Ты надел %s в слот сумки\n> ", item.GetColoredName(it))
	//кольца
	case "ring":
		if p.Equipment.Ring1 == nil {
			p.Equipment.Ring1 = it
			fmt.Fprintf(conn, "Ты надел %s на левую руку\n> ", item.GetColoredName(it))
		} else if p.Equipment.Ring2 == nil {
			p.Equipment.Ring2 = it
			fmt.Fprintf(conn, "Ты надел %s на правую руку\n> ", item.GetColoredName(it))
		} else {
			fmt.Fprintf(conn, "У тебя уже есть два кольца. Сними одно сначало\n> ")
			return
		}
		p.RemoveItemFromStorage(itemName, inBag, index)

	default:
		fmt.Fprintf(conn, "Этот предмет нельзя надеть. Тип : %s\n> ", it.ItemType)
		return
	}
	playerRepo.Save(p)

}
