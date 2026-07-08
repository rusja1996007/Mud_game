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

func HandleDrink(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if cmd == "drink" {
		fmt.Fprintf(conn, "Что выпить?\n> ")
		return
	}

	args, found := strings.CutPrefix(cmd, "drink ")
	if !found {
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		fmt.Fprintf(conn, "Что выпить?\n> ")
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
		fmt.Fprintf(conn, "Предмет не найден\n> ")
		return
	}

	var thatItem *item.ItemStack
	if inBag {
		thatItem = p.Equipment.BagItems[index]
	} else {
		thatItem = p.Inventory[index]
	}

	if thatItem.ItemType != "drink" {
		fmt.Fprintf(conn, "Это нельзя пить\n> ")
		return
	}

	p.Stats.Thirst += thatItem.ThirstRestore
	if p.Stats.Thirst > 100 {
		p.Stats.Thirst = 100
	}

	// Если это вода — возвращаем пустую бутылку
	if thatItem.Name == "water bottle" {
		emptyBottle := item.GetItem("empty bottle", 1)
		if !p.AddItemToInventory(emptyBottle) {
			// Если нет места — кидаем на пол
			room, err := roomRepo.FindByID(p.CurrentRoom)
			if err == nil {
				room.AddItem(emptyBottle)
				roomRepo.Save(room)
				fmt.Fprintf(conn, "Ты выпил %s, бутылка упала на пол. Жажда: %d/100\n> ", itemName, p.Stats.Thirst)

			}

		} else {
			fmt.Fprintf(conn, "Ты выпил %s, бутылку оставил. Жажда: %d/100\n> ", itemName, p.Stats.Thirst)
		}

		// Удаляем 1 единицу напитка
		p.RemoveOneItem(itemName, inBag, index)
		playerRepo.Save(p)
	}
}
