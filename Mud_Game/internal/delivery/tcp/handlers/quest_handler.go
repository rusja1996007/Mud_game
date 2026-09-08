package handlers

import (
	"Mud_game/Mud_Game/internal/domain/npc"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/quest"
	"fmt"
	"net"
)

// обработка ответа на квест
func HandleQuestAnswer(conn net.Conn, cmd string, p *player.Player) bool {
	if !p.PendingQuest {
		return false
	}

	if cmd == "yes" {
		p.ActiveQuests = append(p.ActiveQuests, p.PendingQuestID)
		p.PendingQuest = false
		fmt.Fprintf(conn, "Добавлено задание\n> ")
		return true
	}

	if cmd == "no" {
		p.PendingQuest = false
		fmt.Fprintf(conn, "Вы отказались\n> ")
		return true
	}

	fmt.Fprintf(conn, "Напиши `yes` или `no`\n> ")
	return true
}

// Журнал заданий активных
func HandleQuestList(conn net.Conn, p *player.Player) {
	if len(p.ActiveQuests) == 0 {
		fmt.Fprintf(conn, "Журнал заданий пуст\n> ")
		return
	}

	if len(p.ActiveQuests) > 0 {
		fmt.Fprintf(conn, "📜 Активные квесты:\n")
		for _, id := range p.ActiveQuests {
			q, ok := quest.Quests[id]
			if !ok {
				continue
			}
			fmt.Fprintf(conn, "• %s\n %s\n", q.Name, q.Description)
		}
	}
	fmt.Fprintf(conn, "> ")
}

// квест "кольцо для деда"
func QUESTRingQuest(conn net.Conn, p *player.Player, target *npc.NPC) bool {

	//завершен ли
	q := quest.Quests["ring_quest"]
	for _, id := range p.CompletedQuests {
		if id == "ring_quest" {
			fmt.Fprintf(conn, "🧓 Спасибо тебе ещё раз за кольцо! Я теперь спокоен.\n> ")
			return true
		}
	}

	hasQuest := p.HasActiveQuest("ring_quest")        //активен ли квест
	idx, inBag := p.FindItemGlobalByName("wife ring") //есть ли кольцо

	fmt.Fprintf(conn, "🧓 %s\n", target.Description)
	if !hasQuest {
		fmt.Fprintf(conn, "Дед: Внучек, помоги мне! Гоблин украл кольцо моей жены...\n")
		fmt.Fprintf(conn, "Если вернешь кольцо, я щедро награжу тебя!\n")
		fmt.Fprintf(conn, "Ты получишь:\n +%d EXP \n +%d монет\n", q.ExpReward, q.CoinReward)
		fmt.Fprintf(conn, "Согласен? (напиши 'yes' или 'no')\n> ")
		p.PendingQuest = true
		p.PendingQuestID = "ring_quest"
		return true
	}

	//если активен но кольца нету
	if hasQuest && idx == -1 {
		fmt.Fprintf(conn, "Я бы этих гоблинов в молодости с одним ножом всех обосал...\nно уже возраст не тот...\n> ")
		return true
	}

	//сдача
	if hasQuest && idx != -1 {
		p.RemoveOneItem("wife ring", inBag, idx)
		p.CompleteQuest("ring_quest")
		p.AddExperience(q.ExpReward, conn)
		player.AddItem(&p.Inventory, "coin", q.CoinReward)
		//если захочу добавить лут - использовать этот код:
		/*for _, itemName := range q.ItemReward {
			item := item.GetItem(itemName, 1)
			if item != nil {
				p.AddItemToInventory(item)
			}
		}*/

		fmt.Fprintf(conn, "🧓 Спасибо! Вот твоя награда:\n+%d EXP \n+%d монет\n> ", q.ExpReward, q.CoinReward)
		return true
	}
	return false
}
