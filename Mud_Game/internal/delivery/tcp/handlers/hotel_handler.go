package handlers

import (
	"Mud_game/Mud_Game/internal/domain/buff"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"time"
)

func HandleHotel(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if p.CurrentRoom != "hotel" {
		fmt.Fprintf(conn, "Ты не в гостинице\n> ")
		return
	}

	if p.Stats.IsSleeping {
		fmt.Fprintf(conn, "Ты уже спишь\n> ")
		return
	}

	if !player.HasItem(p.Inventory, "coin", 20) {
		fmt.Fprintf(conn, "Недостаточно монет, надо 20\n> ")
		return
	}

	player.RemoveItem(&p.Inventory, "coin", 20)

	p.Stats.IsSleepingHotel = true
	p.Stats.SleepStartTime = time.Now()
	playerRepo.Save(p)

	fmt.Fprintf(conn, "Ты заплатил 20 монет и забрел в гостиницу. Отдых займет 10 минут.\n> ")

	go func() {
		time.Sleep(5 * time.Second) ///////////////////ВРЕМЕННО
		p.Stats.Health += 60        //восстанавливает 60хп
		maxHealth := 50 + p.Stats.Strength*5
		if p.Stats.Health > maxHealth {
			p.Stats.Health = maxHealth
		}

		p.Stats.Hunger = 100
		p.Stats.Thirst = 100

		newBuff := &buff.Buff{
			ID:            "hotel_health_boost",
			Type:          buff.MaxHealthBoost,
			Value:         50,
			Duration:      20 * time.Second, //.......................ВРЕМЕННО
			RemainingTime: 20 * time.Second, //....................ВРЕМЕННО
			Description:   "+50 к максимальному здоровью",
		}

		//добавление и применение бафа:
		p.ActiveBuffs = append(p.ActiveBuffs, newBuff)
		p.ApplyBuffEffect(newBuff)

		p.Stats.IsSleepingHotel = false
		p.Stats.SleepStartTime = time.Time{}

		playerRepo.Save(p)
		p.SendMessage(conn, "\nТы хорошо отдохнул.\nВосстановил 60 HP\nГолод и жажда восстановлена\n+50 к макс HP на 1 час\n> ")
	}()

}
