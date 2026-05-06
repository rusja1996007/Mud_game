package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
)

func HandleStats(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	maxHP := 50 + p.Stats.Strength*5
	maxEXP := p.GetExpForNextLevel()
	fmt.Fprintf(conn, "======ХАРАКТЕРИСТИКИ======\n")
	fmt.Fprintf(conn, "Здоровье: %d/%d\n", p.Stats.Health, maxHP)
	fmt.Fprintf(conn, "Голод: %d/100\n", p.Stats.Hunger)
	fmt.Fprintf(conn, "Жажда: %d/100\n", p.Stats.Thirst)
	fmt.Fprintf(conn, "==========================\n")
	fmt.Fprintf(conn, "Сила: %d\n", p.Stats.Strength)
	fmt.Fprintf(conn, "Ловкость: %d\n", p.Stats.Dexterity)
	fmt.Fprintf(conn, "Интеллект: %d\n", p.Stats.Intelect)
	fmt.Fprintf(conn, "Следопытство: %d\n", p.Stats.Tracking)
	fmt.Fprintf(conn, "==========================\n")
	fmt.Fprintf(conn, "Уровень: %d\n", p.Stats.Level)
	fmt.Fprintf(conn, "Опыт: %d/%d\n", p.Stats.Experience, maxEXP)
	fmt.Fprintf(conn, "Слоты: %d/%d\n", p.GetUsedSlots(), p.GetMaxSlots())
	fmt.Fprintf(conn, "==========================\n> ")

}
