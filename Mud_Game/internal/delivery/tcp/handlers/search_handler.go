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

	if p.IsSearching {
		fmt.Fprintf(conn, "Ты уже все осматриваешь.Подожди.\n> ")
		return
	}

	if p.CurrentRoom != "dungeon_goblin" && p.CurrentRoom != "dungeon_goblins_v2" && p.CurrentRoom != "glubini_room" {
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

	// ✅ Выбор таблицы лута по комнате
	var lootTable []loot.CaveLootItem
	switch p.CurrentRoom {
	case "dungeon_goblin":
		lootTable = loot.CaveLootTable
	case "dungeon_goblins_v2":
		lootTable = loot.CaveV2LootTable
	case "glubini_room":
		lootTable = loot.GLubiniLootTable
	}

	p.IsSearching = true
	fmt.Fprintf(conn, "Ты начинаешь обыскивать пещеру...\n")

	go func() {
		defer func() {
			p.IsSearching = false
		}()
		time.Sleep(10 * time.Second)

		if (p.CurrentRoom == "dungeon_goblin" || p.CurrentRoom == "glubini_room" || p.CurrentRoom == "dungeon_goblins_v2") && time.Now().Before(monster.TimeToLoot) {

			var foundItems []*item.ItemStack

			//генерируем лут

			found := false

			for _, lootItem := range lootTable {

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
					foundItems = append(foundItems, itemStack)
					found = true
				}

			}
			if found {
				p.SendMessage(conn, "Вы обнаружили:\n")
				for _, item := range foundItems {
					p.SendMessage(conn, fmt.Sprintf(" - %s x%d\n", item.Name, item.Count))
				}
				p.SendMessage(conn, "> ")
			} else {
				p.SendMessage(conn, "Ты обыскал пещеру, но ничего не нашёл.\n> ")
			}
		} else if p.CurrentRoom == "dungeon_goblin" || p.CurrentRoom == "glubini_room" || p.CurrentRoom == "dungeon_goblins_v2" {
			p.SendMessage(conn, "Ты не успел закончить обыск — пещера обвалилась!\n> ")
		}
	}()

	fmt.Fprintf(conn, "> ")
}
