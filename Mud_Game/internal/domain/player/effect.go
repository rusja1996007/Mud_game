package player

import (
	"Mud_game/Mud_Game/internal/domain/buff"
	"fmt"
	"net"
	"time"
)

// применяет эффект от сьеденого/выпитого предмета
func (p *Player) ApplyItemEffect(itemName string, conn net.Conn) {
	switch itemName {
	case "vegetable set":
		p.Stats.Health += 5
		maxHealth := 50 + p.Stats.Strength*5
		if p.Stats.Health > maxHealth {
			p.Stats.Health = maxHealth
		}

		//Баф регенерации
		newBuff := &buff.Buff{
			ID:            "vegetable_set",
			Type:          buff.HealthRegen,
			Value:         1,
			Interval:      30 * time.Second,
			Duration:      10 * time.Minute,
			RemainingTime: 10 * time.Minute,
			Description:   "+1 HP/ 30 сек",
		}
		p.ActiveBuffs = append(p.ActiveBuffs, newBuff)

		fmt.Fprintf(conn, "Вы восстановили 5 HP. Активирована регенерация (+1 HP/30 сек) на 10 минут!\n")

	default:
		return

	}

}
