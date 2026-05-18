package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"time"
)

func HandleWake(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if !p.Stats.IsSleeping {
		fmt.Fprintf(conn, "Персонаж не спит.\n> ")
		return
	}

	proshlo := time.Since(p.Stats.SleepStartTime)
	minutes := int(proshlo.Minutes())

	if minutes > 0 {
		p.Stats.Health += minutes
		maxHealth := 50 + p.Stats.Strength*5
		if p.Stats.Health > maxHealth {
			p.Stats.Health = maxHealth
		}
		fmt.Fprintf(conn, "Ты проснулся, восстановил %d здоровья.\n", minutes)
	} else {
		fmt.Fprintf(conn, "Ты проспал мало, здоровье не восстановилось.\n")
	}

	p.Stats.IsSleeping = false
	p.Stats.SleepStartTime = time.Time{}
	playerRepo.Save(p)
	fmt.Fprintf(conn, ">")

}
