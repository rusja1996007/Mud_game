package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"fmt"
	"net"
)

// получить список всех торговцев в комнате
func (t *TraderHandler) HandleMarket(conn net.Conn, p *player.Player) {
	// торговцы в комнате
	traders, err := t.traderRepo.FindByRoom(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка поиска торговца: %v\n> ", err)
		return
	}

	if len(traders) == 0 {
		fmt.Fprintf(conn, "В комнате нет торговцев\n> ")
		return
	}

	fmt.Fprintf(conn, "🏪 Торговцы на рынке:\n\n")
	for _, trader := range traders {
		fmt.Fprintf(conn, "%s\n", trader.Name)
		fmt.Fprintf(conn, "%s\n", trader.Description)
		fmt.Fprintf(conn, "===================> ")
	}

	fmt.Fprintf(conn, "\nИспользуй talk <имя> чтобы поговорить\n> ")
}
