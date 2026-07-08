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

// сьесть
func HandleEat(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if cmd == "eat" {
		fmt.Fprintf(conn, "Что сьесть?\n> ")
		return
	}

	args, found := strings.CutPrefix(cmd, "eat ")
	if !found {
		return
	}
	args = strings.TrimSpace(args)
	if args == "" {
		fmt.Fprintf(conn, "Что сьесть?\n> ")
		return
	}

	parts := strings.Fields(args)
	var itemName string

	if len(parts) == 1 {
		itemName = parts[0]
	} else {
		itemName = strings.Join(parts, " ")
	}

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
		if num <= len(p.Inventory) {
			inBag = false
		} else {
			inBag = true
		}
	} else {
		index, inBag = p.FindItemGlobalByName(itemName)
	}

	if index == -1 {
		fmt.Fprintf(conn, "Предмет не найден в инвентаре\n> ")
		return
	}

	//получаем сам предмет
	var item *item.ItemStack
	if inBag {
		item = p.Equipment.BagItems[index]
	} else {
		item = p.Inventory[index]
	}

	if item.ItemType != "food" {
		fmt.Fprintf(conn, "Это нельзя есть\n> ")
		return
	}

	//восстанавливаем голод

	p.Stats.Hunger += item.HungerRestore
	if p.Stats.Hunger > 100 {
		p.Stats.Hunger = 100
	}
	p.ApplyItemEffect(itemName, conn)

	//удаляем 1 единицу самой еды
	p.RemoveOneItem(itemName, inBag, index)

	//СОхраняем игрока
	playerRepo.Save(p)

	fmt.Fprintf(conn, "Ты сьел %s. Голод: %d/100\n> ", itemName, p.Stats.Hunger)

}
