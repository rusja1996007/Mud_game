package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
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

	index := p.FindItemIndex(itemName)
	if index == -1 {
		fmt.Fprintf(conn, "Предмет не найден\n> ")
		return
	}

	item := p.Inventory[index]

	if item.ItemType != "drink" {
		fmt.Fprintf(conn, "Это нельзя пить\n> ")
		return
	}

	p.Stats.Thirst += item.ThirstRestore
	if p.Stats.Thirst > 100 {
		p.Stats.Thirst = 100
	}

	player.RemoveItem(&p.Inventory, itemName, 1)

	playerRepo.Save(p)

	fmt.Fprintf(conn, "Ты выпил %s. Жажда:%d/100\n> ", itemName, p.Stats.Thirst)
}
