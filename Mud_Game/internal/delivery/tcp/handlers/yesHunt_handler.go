package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"time"
)

func HandleYesHunt(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if !p.PendingHunt {
		fmt.Fprintf(conn, "Нет ожидающий подтверждений.\n> ")
		return
	}

	//проверяем не истекло ли время
	if time.Now().After(p.PendingHuntExpiry) {
		p.PendingHunt = false
		fmt.Fprintf(conn, "Время подтверждения истекло. Отправь 'hunt' заново.\n> ")
		return
	}

	p.PendingHunt = false
	p.StartHunt(conn, playerRepo, roomRepo)
	fmt.Fprintf(conn, "Ты отправился на охоту! Вернешься через 1 час.\n> ")
}
