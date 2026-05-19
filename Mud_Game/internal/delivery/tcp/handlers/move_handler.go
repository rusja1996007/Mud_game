package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"strings"
	"time"

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
	exits := room.GetExits() //получить карту выходов

	nextRoomID, ok := exits[direction] //Проверить, есть ли такое направление(direction)
	if !ok {
		fmt.Printf("DEBUG: Выход '%s' не найден!\n", direction)
		fmt.Fprintf(conn, "Туда нельзя идти\n> ")
		return
	}

	//===========================Путешествие===============================
	//Если игрок дома и идет домой
	if p.CurrentRoom == p.Zone.RoadID && direction == "south" {
		if p.Stats.Hunger < 10 {
			fmt.Fprintf(conn, "Ты голоден для далеких прогулок.\n> ")
			return
		}

		if p.Stats.Thirst < 20 {
			fmt.Fprintf(conn, "Тебе нужно хорошо попить перед путешествием.\n> ")
			return
		}

		//запрос подтверждения
		p.PendingTravel = true
		p.PendingTravelDirection = direction
		p.PendingTravelExpiry = time.Now().Add(15 * time.Second)

		fmt.Fprintf(conn, "⚠️ Путешествие в город займёт 5 минут реального времени.\n")
		fmt.Fprintf(conn, "Потребуется: 10 голода и 20 жажды.\n")
		fmt.Fprintf(conn, "Ты не сможешь управлять персонажем до окончания.\n")
		fmt.Fprintf(conn, "Напиши 'yes' для подтверждения или 'no' для отмены.\n> ")
		return
	}

	// Если игрок в городе и идёт домой
	if p.CurrentRoom == "global_town" && strings.HasPrefix(direction, "дом ") {
		if p.Stats.Hunger < 10 {
			fmt.Fprintf(conn, "Ты слишком голоден для путешествия. Поешь сначала.\n> ")
			return
		}
		if p.Stats.Thirst < 20 {
			fmt.Fprintf(conn, "Ты слишком хочешь пить для путешествия. Попей сначала.\n> ")
			return
		}

		// Запрашиваем подтверждение
		p.PendingTravel = true
		p.PendingTravelDirection = direction
		p.PendingTravelExpiry = time.Now().Add(15 * time.Second)

		fmt.Fprintf(conn, "⚠️ Путешествие домой займёт 5 минут реального времени.\n")
		fmt.Fprintf(conn, "Потребуется: 10 голода и 20 жажды.\n")
		fmt.Fprintf(conn, "Ты не сможешь управлять персонажем до окончания.\n")
		fmt.Fprintf(conn, "Напиши 'yes' для подтверждения.\n> ")
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
