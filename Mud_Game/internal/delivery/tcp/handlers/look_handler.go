package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
)

func HandleLook(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	currentRoomID := p.CurrentRoom                // 1. Получить текущую комнату игрока(id комнаты)
	room, err := roomRepo.FindByID(currentRoomID) //Обращаемся к репозиторию комнат (s.roomRepo) и просим найти комнату по этому ID.
	if err != nil {
		fmt.Fprintf(conn, "Комната  не найдена\n> ")
		return
	}
	responce := room.Look(p.ID)           //Получаем описание комнаты
	fmt.Fprintf(conn, "%s\n> ", responce) //Добавляем приглашение \n> для следующей команды
}
