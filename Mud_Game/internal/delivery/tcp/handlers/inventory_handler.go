package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func HandleInventory(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//ИНВЕНТАРЬ

	if len(p.Inventory) == 0 {
		fmt.Fprintf(conn, "Инвентарь пуст\n> ")
		return
	}
	// Инвентарь уже в виде стопок, просто выводим!
	//Красивый вывод
	var builder strings.Builder
	builder.WriteString("Твой инвентарь:\n")

	for _, stack := range p.Inventory {
		builder.WriteString(" • ")
		builder.WriteString(stack.Name)
		if stack.Count > 1 {
			builder.WriteString(" x")
			builder.WriteString(strconv.Itoa(stack.Count))
		}
		builder.WriteString("\n")
	}
	builder.WriteString("> ")
	conn.Write([]byte(builder.String()))

}
