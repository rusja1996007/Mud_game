package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
)

func HandleQuit(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//ПОКИНУТЬ ПРИЛОЖЕНИЕ

	fmt.Fprintf(conn, "До свидания!\n")

}
