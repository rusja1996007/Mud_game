package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"time"
)

func HandleSleep(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if p.CurrentRoom != p.Zone.HomeRoomID {
		fmt.Fprintf(conn, "Можешь уснуть только дома\n> ")
		return
	}

	if p.Stats.Hunger < 10 {
		fmt.Fprintf(conn, "Из-за голода ты не сможешь уснуть\n> ")
		return
	}

	if p.Stats.Thirst < 10 {
		fmt.Fprintf(conn, "Из-за сильной жажды ты не сможешь уснуть\n> ")
		return
	}

	if p.Stats.IsSleeping == true {
		fmt.Fprintf(conn, "Ты уже спишь! Проснись командой 'wake'.\n> ")
		return
	}

	p.Stats.IsSleeping = true
	p.Stats.SleepStartTime = time.Now()

	playerRepo.Save(p)

	fmt.Fprintf(conn, "Вы легли спать, жизни восстанавливаются 1HP/минуту.\n> ")
}
