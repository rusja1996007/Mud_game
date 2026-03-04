package handlers

import (
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
	// Ищем нужную стопку в инвентаре
	foundIndex := -1
	for i, stack := range p.Inventory {
		if stack.Name == itemName {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		fmt.Fprintf(conn, "У тебя нет такого предмета\n> ")
		return
	}
	//в переменную ложим количество которое имеем
	available := p.Inventory[foundIndex].Count
	//Определить, сколько уничтожать
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
	//    Уничтожить предметы (удалить из инвентаря)
	if destroyCount == available {
		// Удаляем всю стопку
		p.Inventory = append(p.Inventory[:foundIndex], p.Inventory[foundIndex+1:]...)
	} else {
		// Уменьшаем количество
		p.Inventory[foundIndex].Count -= destroyCount
	}
	playerRepo.Save(p)
	fmt.Fprintf(conn, "Ты уничтожил %d %s\n> ", destroyCount, itemName) //fmt.Fprintf(conn, ...),Пишет напрямую в соединение(ненадо отделдьно преобразовывать в байты) короче и эффективнее

}
