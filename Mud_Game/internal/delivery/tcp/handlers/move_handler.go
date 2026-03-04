package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"strings"

	"net"
)

func HandleMove(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if cmd == "move" {
		fmt.Fprintf(conn, "Куда идти?\n> ")
		return
	}
	direction := strings.TrimPrefix(cmd, "move ") //убирает "move " и возвращает направление
	direction = strings.TrimSpace(direction)
	if direction == "" {
		fmt.Fprintf(conn, "Куда идти?\n> ")
		return
	}
	room, err := roomRepo.FindByID(p.CurrentRoom) //Обращаемся к репозиторию комнат (s.roomRepo) и просим найти комнату по этому ID.
	if err != nil {
		fmt.Fprintf(conn, "Комната  не найдена\n> ")
		return
	}
	exits := room.GetExits()           //получить карту выходов
	nextRoomID, ok := exits[direction] //Проверить, есть ли такое направление(direction)
	if !ok {
		fmt.Fprintf(conn, "Туда нельзя идти\n> ")
		return
	}
	p.CurrentRoom = nextRoomID //Обновить позицию игрока и сохранить
	playerRepo.Save(p)

	nextRoom, _ := roomRepo.FindByID(nextRoomID) //Показываем новую комнату
	fmt.Fprintf(conn, "%s\n> ", nextRoom.Look(p.ID))

}
