package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
)

func HandleFill(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if p.CurrentRoom != p.Zone.WellID {
		fmt.Fprintf(conn, "Ты не у колодца\n> ")
		return
	}

	index := p.FindItemIndex("empty bottle")
	if index == -1 {
		fmt.Fprintf(conn, "У тебя нету пустой бутылки\n> ")
		return
	}

	emptyStack := p.Inventory[index]

	//удаляем пустую бутылку
	if emptyStack.Count > 1 {
		emptyStack.Count--
	} else {
		p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	}

	//если нет места - возвращаем бутылку
	if !p.CanAddItem() {
		p.AddItemToInventory(item.GetItem("empty bottle", 1))
		fmt.Fprintf(conn, "Нет места для полной бутылки!\n> ")
		return
	}

	p.AddItemToInventory(item.GetItem("water bottle", 1))

	playerRepo.Save(p)
	fmt.Fprintf(conn, "Ты наполнил бутылку\n> ")
}
