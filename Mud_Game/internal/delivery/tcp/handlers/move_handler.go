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

	//"Смотрим, есть ли в названии комнаты слова home_, garden_ или road_"

	//Если есть — значит это личная зона (дом, огород или дорога).
	//Если нет — значит это общая комната (город) и можно пускать без проверки.

	if strings.Contains(nextRoomID, "home_") ||
		strings.Contains(nextRoomID, "garden_") ||
		strings.Contains(nextRoomID, "road_") {
		//"Смотрим, есть ли в названии комнаты ID текущего игрока"
		if !strings.Contains(nextRoomID, p.ID) {
			fmt.Fprintf(conn, "Чужая территория, туда нельзя\n> ")
			return
		}
	}
	p.CurrentRoom = nextRoomID //Обновить позицию игрока и сохранить
	playerRepo.Save(p)

	nextRoom, _ := roomRepo.FindByID(nextRoomID) //Показываем новую комнату
	fmt.Fprintf(conn, "%s\n> ", nextRoom.Look(p.ID))

}
