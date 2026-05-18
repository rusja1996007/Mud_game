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
	fmt.Fprintf(conn, "🛡️ Защита:\n")
	fmt.Fprintf(conn, "Физическая: %d\n", p.GetTotalDefence())
	fmt.Fprintf(conn, "Магическая: %d\n", p.GetTotalMagicDefence())
	fmt.Fprintf(conn, "От яда: %d\n", p.GetTotalPoisonDefence())
	fmt.Fprintf(conn, "От огня: %d\n", p.GetTotalFireDefence())
	fmt.Fprintf(conn, "==========================\n")
	fmt.Fprintf(conn, "Уровень: %d\n", p.Stats.Level)
	fmt.Fprintf(conn, "Опыт: %d/%d\n", p.Stats.Experience, maxEXP)
	fmt.Fprintf(conn, "Слоты: %d/%d\n", p.GetUsedSlots(), p.GetMaxSlots())
	fmt.Fprintf(conn, "==========================\n")

	if len(p.ActiveBuffs) > 0 {
		fmt.Fprintf(conn, "\n✨ Активные эффекты:\n")
		for _, b := range p.ActiveBuffs {
			minutes := int(b.RemainingTime.Minutes())
			seconds := int(b.RemainingTime.Seconds()) % 60
			fmt.Fprintf(conn, "  • %s\n", b.Description)
			fmt.Fprintf(conn, "    осталось: %dм %dс\n", minutes, seconds)
		}
	}
	fmt.Fprintf(conn, ">")

}
