package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/loot"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// обыск локации
func HandleSearch(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	if p.CurrentRoom != "dungeon_goblin" {
		fmt.Fprintf(conn, "Тебе тут нечего обыскивать.\n> ")
		return
	}

	//находим комнату и монстра
	room, _ := roomRepo.FindByID(p.CurrentRoom)
	monster := room.GetMonster()

	if monster == nil || monster.IsAlive {
		fmt.Fprintf(conn, "Сначало победи монстра.\n> ")
		return
	}

	if time.Now().After(monster.TimeToLoot) {
		fmt.Fprintf(conn, "Пещера уже обвалилась, обыскивать поздно.\n> ")
		return
	}

	found := false

	for _, lootItem := range loot.CaveLootTable {

		// шанс
		chance := lootItem.BaseChance + 5*p.Stats.Tracking

		//если успех:
		if rand.Intn(100) < chance {
			//расчет кол-ва
			count := lootItem.MinCount
			if lootItem.MaxCount > lootItem.MinCount {
				count += rand.Intn(lootItem.MaxCount - lootItem.MinCount + 1)
			}
			itemStack := &item.ItemStack{
				Name:     lootItem.ItemData.Name,
				Count:    count,
				ItemType: lootItem.ItemData.ItemType,
			}
			room.AddItem(itemStack)
			roomRepo.Save(room)
			fmt.Fprintf(conn, "В пещере ты нашел %s x%d \n", itemStack.Name, itemStack.Count)
			found = true
		}

	}
	if !found {
		fmt.Fprintf(conn, "Ты обыскал пещеру, но ничего не нашел.\n")
	}
	fmt.Fprintf(conn, "> ")
}
