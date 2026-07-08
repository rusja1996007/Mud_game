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

func HandleDestroy(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//УНИЧТОЖИТЬ destroy
	if cmd == "destroy" {
		fmt.Fprintf(conn, "Что уничтожить?\n> ")
		return
	}
	argss, found := strings.CutPrefix(cmd, "destroy ")
	if !found {
		return // это не destroy - идем дальше
	}
	argss = strings.TrimSpace(argss)
	if argss == "" {
		fmt.Fprintf(conn, "Что уничтожить?\n> ")
		return
	}
	parts := strings.Fields(argss) //парсим(разбиваем) команду
	var count int = 1
	var itemName string

	if len(parts) == 1 {
		itemName = parts[0]
	} else if len(parts) >= 2 {
		if parts[0] == "all" {
			count = -1
			itemName = strings.Join(parts[1:], " ")
		} else {
			num, err := strconv.Atoi(parts[0])
			if err == nil && num > 0 {
				count = num
				itemName = strings.Join(parts[1:], " ")
			} else if err == nil && num <= 0 {
				fmt.Fprintf(conn, "Количество должно быть положительным\n> ")
				return
			} else {
				itemName = strings.Join(parts, " ")
			}

		}
	}

	if itemName == "" {
		fmt.Fprintf(conn, "Что уничтожить?\n> ")
		return
	}
	// Ищем предмет в инвентаре и мешке
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
		fmt.Fprintf(conn, "У тебя нет такого предмета\n> ")
		return
	}

	// Получаем предмет
	var thatItem *item.ItemStack
	if inBag {
		thatItem = p.Equipment.BagItems[index]
	} else {
		thatItem = p.Inventory[index]
	}

	available := thatItem.Count
	destroyCount := count
	if count == -1 {
		destroyCount = available
	}
	if destroyCount > available {
		destroyCount = available
	}
	if destroyCount == 0 {
		fmt.Fprintf(conn, "Нечего уничтожать\n> ")
		return
	}

	// Удаляем
	if destroyCount == available {
		p.RemoveItemFromStorage(itemName, inBag, index)
	} else {
		if inBag {
			p.Equipment.BagItems[index].Count -= destroyCount
		} else {
			p.Inventory[index].Count -= destroyCount
		}
	}

	playerRepo.Save(p)
	fmt.Fprintf(conn, "Ты уничтожил %d %s\n> ", destroyCount, itemName)
}
