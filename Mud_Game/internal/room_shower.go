package interfaces

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"net"
)

// RoomShower - интерфейс для показа комнаты с NPC
type RoomShower interface {
	ShowRoomWithNPC(conn net.Conn, p *player.Player)
}
