package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleHunt(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	if p.CurrentRoom != p.Zone.HomeRoomID {
		fmt.Fprintf(conn, "Ты должен быть в доме чтобы пойти на охоту!\n> ")
		return
	}

	waterIndex := p.FindItemIndex("water bottle")
	if waterIndex == -1 || p.Inventory[waterIndex].Count < 2 {
		fmt.Fprintf(conn, "Нужны 2 бутылки с водой для охоты!\n> ")
		return
	}

	if p.Stats.IsHunting {
		fmt.Fprintf(conn, "Ты на охоте, вернешься через %v\n> ", time.Until(p.Stats.HuntingEndTime).Round(time.Second))
		return
	}

	// Проверяем, есть ли ожидание подтверждения
	if p.PendingHunt {
		if strings.ToLower(cmd) == "yes" {
			//запускаем охоту
			p.PendingHunt = false
			p.StartHunt(conn, playerRepo, roomRepo)
			fmt.Fprintf(conn, "Ты отправился на охоту! Вернешься через 1 час.\n> ")
		} else {
			p.PendingHunt = false
			fmt.Fprintf(conn, "Охота отменена.\n> ")
		}
		return
	}

	//запрашиваем подтверждение
	p.PendingHunt = true
	fmt.Fprintf(conn, "ВНИМАНИЕ! Охота займет 1 час реального времени.\n> ")
	fmt.Fprintf(conn, "Ты не сможешь управлять персонажем до ее окончания.\n> ")
	fmt.Fprintf(conn, "Напиши 'yes' для подтверждения\n> ")
}
