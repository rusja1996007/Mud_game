package handlers

import (
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

	if num, err := strconv.Atoi(itemName); err == nil {
		target, idx := p.FindItemByNumber(num)
		if target == nil {
			fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", num)
			return
		}
		index = idx
		itemName = target.Name
	} else {
		index = p.FindItemIndex(itemName)
	}

	if index == -1 {
		fmt.Fprintf(conn, "У тебя нет предмета %s в инвентаре\n> ", itemName)
		return
	}
	//Получи предмет из инвентаря
	item := p.Inventory[index]

	switch item.ItemType {

	//оружие
	case "weapon":
		if p.Equipment.Weapon != nil {
			p.AddItemToInventory(p.Equipment.Weapon)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Weapon.Name)
		}
		//вооружаем
		p.Equipment.Weapon = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот оружия\n> ", itemName)
	//броня
	case "armor":
		if p.Equipment.Armor != nil {
			p.AddItemToInventory(p.Equipment.Armor)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Armor.Name)
		}
		//вооружаем
		p.Equipment.Armor = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот брони\n> ", itemName)
	//шлем
	case "helmet":
		if p.Equipment.Helmet != nil {
			p.AddItemToInventory(p.Equipment.Helmet)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Helmet.Name)
		}
		//вооружаем
		p.Equipment.Helmet = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот шлема\n> ", itemName)
	//ботинки
	case "boots":
		if p.Equipment.Boots != nil {
			p.AddItemToInventory(p.Equipment.Boots)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Boots.Name)
		}
		//вооружаем
		p.Equipment.Boots = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот обуви\n> ", itemName)
	//щит
	case "shield":
		if p.Equipment.Shield != nil {
			p.AddItemToInventory(p.Equipment.Shield)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Shield.Name)
		}
		//вооружаем
		p.Equipment.Shield = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот щита\n> ", itemName)
	//сумка
	case "bag":
		if p.Equipment.Bag != nil {
			p.AddItemToInventory(p.Equipment.Bag)
			fmt.Fprintf(conn, "Ты снял %s\n", p.Equipment.Bag.Name)
		}
		//вооружаем
		p.Equipment.Bag = item
		//удаляем из инвентаря
		player.RemoveItem(&p.Inventory, itemName, 1)
		fmt.Fprintf(conn, "Ты надел %s в слот сумки\n> ", itemName)
	//кольца
	case "ring":
		if p.Equipment.Ring1 == nil {
			p.Equipment.Ring1 = item
			fmt.Fprintf(conn, "Ты надел %s на левую руку\n> ", itemName)
		} else if p.Equipment.Ring2 == nil {
			p.Equipment.Ring2 = item
			fmt.Fprintf(conn, "Ты надел %s на правую руку\n> ", itemName)
		} else {
			fmt.Fprintf(conn, "У тебя уже есть два кольца. Сними одно сначало\n> ")
			return
		}
		player.RemoveItem(&p.Inventory, itemName, 1)

	//Еда, напитки, семена, материалы
	case "food", "drink", "seed", "material", "container":
		fmt.Fprintf(conn, "Этот предмет нельзя надеть.\n> ")
		return

	default:
		fmt.Fprintf(conn, "Этот предмет нельзя надеть. Тип : %s\n> ", item.ItemType)
		return
	}
	playerRepo.Save(p)

}
