package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
)

// побег
func HandleEscape(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	//текущая комната
	roomInterface, _ := roomRepo.FindByID(p.CurrentRoom)
	concreteRoom, ok := roomInterface.(*room.Room)
	if !ok {
		fmt.Fprintf(conn, "Ошибка приведения комнаты\n> ")
		return
	}

	monster := concreteRoom.GetMonster()

	if monster == nil || !monster.IsAlive {
		fmt.Fprintf(conn, "Некого бояться, можно не убегать.\n> ")
		return
	}

	chanceOfEscape := 50 + 5*p.Stats.Tracking

	//шанс потерять чтото из инвентаря при побеге
	if len(p.Inventory) > 0 {
		//⚠️ Важно: Идём с конца, чтобы не сбивать индексы при удалении.
		for i := len(p.Inventory) - 1; i >= 0; i-- {
			if rand.Intn(100) < 30 {
				lostItem := p.Inventory[i]
				p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
				fmt.Fprintf(conn, "При побеге ты потерял %s\n", lostItem.Name)
			}
		}
	}
	/////////////////////////////////////если успех/////////////////////////////////
	if rand.Intn(100) < chanceOfEscape {
		p.CurrentRoom = concreteRoom.GetExitRoomID()

		// ✅ ВОССТАНАВЛИВАЕМ ВСЕХ МОНСТРОВ
		if len(concreteRoom.MonsterS) > 0 {
			for _, m := range concreteRoom.MonsterS {
				m.Health = m.MaxHealth
				m.IsAlive = true
				if m.ID == "goblin_shaman" {
					m.CastTime = 0
					m.IsCasting = true
				}
			}
		} else if monster != nil {
			monster.Health = monster.MaxHealth
			monster.IsAlive = true
		}

		concreteRoom.ClearItems()
		concreteRoom.SetPlayerOccupantID("")
		playerRepo.Save(p)
		roomRepo.Save(concreteRoom)

		newRoom, _ := roomRepo.FindByID(p.CurrentRoom)
		fmt.Fprintf(conn, "Ты успешно сбежал!\n")
		fmt.Fprintf(conn, "%s\n> ", newRoom.Look(p.ID))
		return
	} else {
		/////////////////////////////////если нет../////////////////////////////////
		defence := p.GetTotalDefence()
		damage := monster.MinDamage + rand.Intn(monster.MaxDamage-monster.MinDamage+1)

		totalDamage := damage - defence

		if totalDamage < 1 {
			totalDamage = 1
		}
		p.Stats.Health -= totalDamage

		fmt.Fprintf(conn, "Ты пытаешься убежать, но монстр успел нанести %d урона вслед!\n", totalDamage)

		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб...\n")
			playerRepo.Delete(p.ID)
			conn.Close()
			return
		}

		p.CurrentRoom = concreteRoom.GetExitRoomID()

		newRoom, _ := roomRepo.FindByID(p.CurrentRoom)

		// ✅ ВОССТАНАВЛИВАЕМ ВСЕХ МОНСТРОВ
		if len(concreteRoom.MonsterS) > 0 {
			for _, m := range concreteRoom.MonsterS {
				m.Health = m.MaxHealth
				m.IsAlive = true
				if m.ID == "goblin_shaman" {
					m.CastTime = 0
					m.IsCasting = true
				}
			}
		} else if monster != nil {
			monster.Health = monster.MaxHealth
			monster.IsAlive = true
		}

		concreteRoom.ClearItems()
		concreteRoom.SetPlayerOccupantID("")
		roomRepo.Save(concreteRoom)
		playerRepo.Save(p)
		fmt.Fprintf(conn, "Ты сбежал!\n")
		fmt.Fprintf(conn, "%s\n> ", newRoom.Look(p.ID))

	}
}
