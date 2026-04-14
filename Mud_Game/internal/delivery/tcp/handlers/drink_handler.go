package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
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

	index := p.FindItemIndex(itemName)
	if index == -1 {
		fmt.Fprintf(conn, "Предмет не найден\n> ")
		return
	}

	thatItem := p.Inventory[index]

	if thatItem.ItemType != "drink" {
		fmt.Fprintf(conn, "Это нельзя пить\n> ")
		return
	}

	p.Stats.Thirst += thatItem.ThirstRestore
	if p.Stats.Thirst > 100 {
		p.Stats.Thirst = 100
	}

	if thatItem.Name == "water bottle" {

		//генерация числа с 0 до 100
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		chance := rng.Intn(100) + 1

		if chance <= 70 {
			p.AddItemToInventory(item.GetItem("empty bottle", 1))
			fmt.Fprintf(conn, "Ты выпил воду, бутылка целая. Жажда:%d/100\n> ", p.Stats.Thirst)
			player.RemoveItem(&p.Inventory, itemName, 1)
			playerRepo.Save(p)
			return
		} else {
			fmt.Fprintf(conn, "Ты выпил воду, но бутылка износилась. Жажда:%d/100\n> ", p.Stats.Thirst)
			player.RemoveItem(&p.Inventory, itemName, 1)
			playerRepo.Save(p)
			return
		}
	}
}
