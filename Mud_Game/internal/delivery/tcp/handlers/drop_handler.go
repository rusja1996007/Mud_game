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

func HandleDrop(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//DROP
	if cmd == "drop" {
		fmt.Fprintf(conn, "Что бросить?\n> ")
		return
	}
	argsss, found := strings.CutPrefix(cmd, "drop ")
	if !found {
		return
	}
	argsss = strings.TrimSpace(argsss)
	if argsss == "" {
		fmt.Fprintf(conn, "Что бросить?\n> ")
		return
	}
	parts := strings.Fields(argsss)
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
				fmt.Fprintf(conn, "Количество должно быть положительныйм\n> ")
				return

			} else {
				itemName = strings.Join(parts, " ")
			}
		}
	}
	if itemName == "" && count == -1 {
		fmt.Fprintf(conn, "Что именно бросить?")
		return
	}

	foundIndex := -1
	for i, stack := range p.Inventory {
		if stack.Name == itemName {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		fmt.Fprintf(conn, "У вас нету такого предмета\n> ")
		return
	}

	available := p.Inventory[foundIndex].Count
	dropCount := count //сколько выложить
	if dropCount == -1 {
		dropCount = available
	}
	if dropCount > available {
		dropCount = available
	}
	if dropCount == 0 {
		fmt.Fprintf(conn, "Нечего бросать\n> ")
		return
	}

	//удаляем из инвентаря игрока
	if dropCount == available {
		p.Inventory = append(p.Inventory[:foundIndex], p.Inventory[foundIndex+1:]...)
	} else {
		p.Inventory[foundIndex].Count -= dropCount
	}

	// ✅ СОЗДАЕМ СТОПКУ ДЛЯ БРОСКА
	dropStack := &item.ItemStack{
		Name:  itemName,
		Count: dropCount,
	}

	//Добавить предметы в комнату
	room, err := roomRepo.FindByID(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты \n> ")
		return
	}

	err = room.AddItem(dropStack)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка при добавлении предмета в комнату\n> ")
		return
	}
	playerRepo.Save(p)
	roomRepo.Save(room)

	fmt.Fprintf(conn, "Ты бросил %d %s\n> ", dropCount, itemName)

}
