package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
)

func HandleStats(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	fmt.Fprintf(conn, "======ХАРАКТЕРИСТИКИ======\n")
	fmt.Fprintf(conn, "Здоровье: %d\n", p.Stats.Health)
	fmt.Fprintf(conn, "Голод: %d\n", p.Stats.Hunger)
	fmt.Fprintf(conn, "Жажда: %d\n", p.Stats.Thirst)
	fmt.Fprintf(conn, "==========================\n")
	fmt.Fprintf(conn, "Сила: %d\n", p.Stats.Strength)
	fmt.Fprintf(conn, "Ловкость: %d\n", p.Stats.Dexterity)
	fmt.Fprintf(conn, "Интеллект: %d\n", p.Stats.Intelect)
	fmt.Fprintf(conn, "Следопытство: %d\n", p.Stats.Tracking)
	fmt.Fprintf(conn, "==========================\n")
	fmt.Fprintf(conn, "Уровень: %d\n", p.Stats.Level)
	fmt.Fprintf(conn, "Опыт: %d\n", p.Stats.Experience)
	fmt.Fprintf(conn, "Слоты: %d/%d\n", p.GetUsedSlots(), p.GetMaxSlots())
	fmt.Fprintf(conn, "==========================\n> ")

}
