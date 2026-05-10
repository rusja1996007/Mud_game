package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"time"
)

// показ и распределдения очков характеристик
func HandleStatPoints(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if p.Stats.PendingStatPoint <= 0 {
		fmt.Fprintf(conn, "У вас нет неиспользованых очков характеристик!\n> ")
		return
	}

	//20 сек на выбор
	p.PendingStatChoiсe = true
	p.PendingStatChoiсeExpiry = time.Now().Add(20 * time.Second)

	fmt.Fprintf(conn, "Количество неиспользуемых очков - %d\n", p.Stats.PendingStatPoint)
	fmt.Fprintf(conn, "Выберите характеристику для повышения:\n")
	fmt.Fprintf(conn, "1.Сила - %d\n", p.Stats.Strength)
	fmt.Fprintf(conn, "2.Ловкость - %d\n", p.Stats.Dexterity)
	fmt.Fprintf(conn, "3.Интелект - %d\n", p.Stats.Intelect)
	fmt.Fprintf(conn, "4.Следопытство - %d\n", p.Stats.Tracking)
	fmt.Fprintf(conn, "Введите номер(1-4):\n> ")
}
