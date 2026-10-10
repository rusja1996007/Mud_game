package common

import (
	interfaces "Mud_game/Mud_Game/internal"
	"Mud_game/Mud_Game/internal/delivery/tcp/handlers"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"net"
)

func RouteCommand(
	conn net.Conn,
	cmd string,
	p *player.Player,
	playerRepo player.Repository,
	roomRepo room.Repository,
	npcRepo *npc_repo.PostgresNPCRepository,
	shower interfaces.RoomShower,
) bool {

	//всегда можем выйти
	if cmd == "quit" {

		//если выход производится во время боя
		if p.CurrentRoom == "dungeon_goblin" {
			room, err := roomRepo.FindByID(p.CurrentRoom)
			if err != nil {
				fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
				return false
			}
			monster := room.GetMonster()

			//если монстр жив-автоматический побег с получением урона
			if monster != nil && monster.IsAlive {
				monsterDamage := monster.MinDamage + rand.Intn(monster.MaxDamage-monster.MinDamage+1)
				defence := p.GetTotalDefence()
				reduction := float64(defence) / (float64(defence) + 100)
				finalDamage := int(float64(monsterDamage) * (1 - reduction))
				if finalDamage <= 0 {
					finalDamage = 1
				}

				p.Stats.Health -= finalDamage
				monster.Health = monster.MaxHealth

				msg := fmt.Sprintf("💨 Перед выходом ты сбегаешь и монстр нанёс %d урона вслед.\n> ", finalDamage)
				p.SendMessage(conn, msg)

				if p.Stats.Health <= 0 {
					p.SendMessage(conn, "💀Ты погиб...\n")
					monster.Health = monster.MaxHealth
					room.SetPlayerOccupantID("")
					roomRepo.Save(room)
					p.StopAllTickers()
					playerRepo.Delete(p.ID)
					conn.Close()
					return true

				}
				//телепорт
				p.CurrentRoom = room.GetExitRoomID()
				room.SetPlayerOccupantID("")
				p.Stats.IsInDungeon = false
				p.Stats.EnteredDungeonAt = time.Time{}
				roomRepo.Save(room)
				playerRepo.Save(p)

				handlers.HandleQuit(conn, cmd, p, roomRepo, playerRepo)
				return true
			}
		}
	}
	//ответ на квест(yes/no)
	if p.PendingQuest {
		handlers.HandleQuestAnswer(conn, cmd, p)
		return false
	}

	//путешествие
	if p.PendingTravel {
		if cmd == "yes" {
			p.PendingTravel = false

			if p.PendingTravelDirection == "south" {
				p.Stats.Hunger -= 10
				p.Stats.Thirst -= 20
				p.Stats.TravelTargetRoom = "global_town"
				fmt.Fprintf(conn, "Ты отправляешься в город. Путь займёт 5 минут.\n> ")
			} else if strings.HasPrefix(p.PendingTravelDirection, "дом ") {
				if p.Zone == nil {
					return false
				}
				p.Stats.Hunger -= 10
				p.Stats.Thirst -= 20

				p.Stats.TravelTargetRoom = p.Zone.RoadID
				fmt.Fprintf(conn, "Ты отправляешься домой. Путь займёт 5 минут.\n> ")
			} else if p.PendingTravelDirection == "dungeon" {
				p.Stats.Hunger -= 5
				p.Stats.Thirst -= 5
				p.Stats.TravelTargetRoom = "dungeon_entrance_goblins"
				fmt.Fprintf(conn, "Ты отправляешься к подземелью. Путь займёт 2 минуты.\n> ")
			}

			p.Stats.IsTraveling = true
			p.Stats.TravelEndTime = time.Now().Add(5 * time.Second) //////////////временно
			playerRepo.Save(p)

			//запуск горутины для автоматического завершения
			go func() {
				time.Sleep(time.Until(p.Stats.TravelEndTime))

				//безопасно отправляем сообщение
				p.SendMessage(conn, "\nТы прибыл!\n")

				//обновляем состояние
				p.CurrentRoom = p.Stats.TravelTargetRoom
				p.Stats.IsTraveling = false
				p.Stats.TravelEndTime = time.Time{}
				p.Stats.TravelTargetRoom = ""
				playerRepo.Save(p)

				//room, _ := s.roomRepo.FindByID(p.CurrentRoom)
				//p.SendMessage(conn, room.Look(p.ID)+"\n> ")
				shower.ShowRoomWithNPC(conn, p)

			}()

			return false
		} else if cmd == "no" {
			p.PendingTravel = false
			fmt.Fprintf(conn, "Путь отменен\n> ")
			return false
		} else {
			fmt.Fprintf(conn, "Сначала подтверди путешествие командой 'yes' или отмени его командой 'no'\n> ")
			return false
		}
	}

	// Если в путешествии — блокируем все команды
	if p.Stats.IsTraveling {
		if time.Now().After(p.Stats.TravelEndTime) {
			p.CurrentRoom = p.Stats.TravelTargetRoom
			p.Stats.IsTraveling = false
			p.Stats.TravelEndTime = time.Time{}
			p.Stats.TravelTargetRoom = ""
			playerRepo.Save(p)

			//room, err := s.roomRepo.FindByID(p.CurrentRoom)
			//if err != nil {
			//	fmt.Fprintf(conn, "Ошибка загрузки комнаты.\n> ")
			//	return false
			//}

			fmt.Fprintf(conn, "Ты прибыл!\n")
			//fmt.Fprintf(conn, "%s\n> ", room.Look(p.ID))
			shower.ShowRoomWithNPC(conn, p)
			return false

		} else {

			remaining := time.Until(p.Stats.TravelEndTime).Round(time.Second)
			fmt.Fprintf(conn, "Ты в пути. Осталось: %v. Команды недоступны.\n> ", remaining)
			return false
		}
	}
	//если спишь, блокируем все
	if p.Stats.IsSleeping && cmd != "wake" {
		fmt.Fprintf(conn, "Ты спишь, проснись командой 'wake'.\n> ")
		return false
	}

	//если в отеле то ждем
	if p.Stats.IsSleepingHotel {
		fmt.Fprintf(conn, "Ты отдыхаешь в отеле. Подожди окончания отдыха.\n> ")
		return false
	}

	// Если игрок на охоте — блокируем все команды кроме "hunt"
	if p.Stats.IsHunting {
		if cmd == "hunt" {
			fmt.Fprintf(conn, "Ты на охоте, вернешься через %v\n> ",
				time.Until(p.Stats.HuntingEndTime).Round(time.Second))

		} else if cmd == "quit" {
			handlers.HandleQuit(conn, cmd, p, roomRepo, playerRepo)
		} else {
			fmt.Fprintf(conn, "Ты на охоте! Нельзя использовать команды кроме hunt и quit\n> ")
		}
		return false
	}

	//Если разговаривает - блокирует все кроме stop talk и quit
	if p.IsTalkin && isCommandBlockedInDialog(cmd, p) {
		fmt.Fprintf(conn, "Ты не можешь использовать команды во время разговора! Используй 'stop talk'.\n> ")
		return false
	}

	//обработка выбора характеристик
	if p.PendingStatChoiсe {
		if time.Now().After(p.PendingStatChoiсeExpiry) {
			p.PendingStatChoiсe = false
			fmt.Fprintf(conn, "Время ожидания истекло, введите повторно 'statpoints'\n> ")
			return false
		}

		if cmd == "1" || cmd == "2" || cmd == "3" || cmd == "4" {
			p.ProcessStatChoice(cmd, conn)
			return false
		} else {
			fmt.Fprintf(conn, "Некорректный ввод. Введите 1,2,3 или 4\n> ")
			return false
		}
	}

	if p.PendingHunt && cmd != "yes" {
		p.PendingHunt = false
		fmt.Fprintf(conn, "Подтверждение охоты отменено.\n> ")
		return false
	}

	if cmd == "" {
		handlers.HandleEmpty(conn, cmd, p, roomRepo, playerRepo)
		return false
	}

	switch {
	///////////////////////////////////////////////////////////////////////////////////
	case cmd == "damage":
		p.Stats.Health -= 40
		if p.Stats.Health < 1 {
			p.Stats.Health = 1
		}
		fmt.Fprintf(conn, "Здоровье уменьшено до %d\n> ", p.Stats.Health)
		return false
		//////////////////////////////////////////////////////////////////////
	case cmd == "sleep":
		handlers.HandleSleep(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "wake":
		handlers.HandleWake(conn, cmd, p, roomRepo, playerRepo)
		return false
	case cmd == "yes":
		handlers.HandleYesHunt(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "hunt":
		handlers.HandleHunt(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "quit":

		handlers.HandleQuit(conn, cmd, p, roomRepo, playerRepo)
		return true // сигнал на выход

	case cmd == "inventory":
		handlers.HandleInventory(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "stats":
		handlers.HandleStats(conn, cmd, p, roomRepo, playerRepo)
		return false
	case cmd == "statpoints":
		handlers.HandleStatPoints(conn, cmd, p, roomRepo, playerRepo)
		return false
	case cmd == "all npc":
		handlers.HandleAllNPC(conn, npcRepo, p)
		return false
	case cmd == "stop talk":
		handlers.HandleStopTalk(conn, p)
		shower.ShowRoomWithNPC(conn, p)
		return false
	case cmd == "look":
		handlers.HandleLook(conn, cmd, p, roomRepo, playerRepo)
		shower.ShowRoomWithNPC(conn, p)
		return false

	case strings.HasPrefix(cmd, "look "):
		handlers.HandleLook(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "move", strings.HasPrefix(cmd, "move "): //после move идет еще чтото. аналогично ниже
		handlers.HandleMove(conn, cmd, p, roomRepo, playerRepo, shower)
		return false

	case strings.HasPrefix(cmd, "take "):
		//если игрок разговаривает с кузнецом -take для кузнеца
		if p.IsTalkin && p.TalkToID == "blacksmith" {
			handlers.HandleTakeRepaired(cmd, conn, p, npcRepo)
			return false
		}
		//иначе - обычный take
		handlers.HandleTake(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "repair "):
		handlers.HandleRepair(conn, cmd, p, npcRepo)
		return false

	case strings.HasPrefix(cmd, "drop "):

		handlers.HandleDrop(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "destroy "):

		handlers.HandleDestroy(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "garden"):
		handlers.HandleGarden(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "plant "):
		handlers.HandlePlant(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "harvest "):
		handlers.HandleHarvest(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "wear "):
		handlers.HandleWear(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "remove "):
		handlers.HandleRemove(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "eat "):
		handlers.HandleEat(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "drink "):
		handlers.HandleDrink(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "fill "):
		handlers.HandleFill(conn, cmd, p, roomRepo, playerRepo)
		return false

	case cmd == "craft":
		handlers.HandleCraft(conn, cmd, p, roomRepo, playerRepo)
		return false
	case strings.HasPrefix(cmd, "craft "):
		handlers.HandleCraft(conn, cmd, p, roomRepo, playerRepo)
		return false
	case strings.HasPrefix(cmd, "use "):
		handlers.HandleUse(conn, cmd, p, roomRepo, playerRepo)
		return false
	case strings.HasPrefix(cmd, "pay 20"):
		handlers.HandleHotel(conn, cmd, p, roomRepo, playerRepo)
		shower.ShowRoomWithNPC(conn, p)
		return false

	case strings.HasPrefix(cmd, "attack"):
		handlers.HandleAttack(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "search"):
		handlers.HandleSearch(conn, cmd, p, roomRepo, playerRepo)
		return false

	case strings.HasPrefix(cmd, "escape"):
		handlers.HandleEscape(conn, cmd, p, roomRepo, playerRepo)
		return false
	case cmd == "talk", strings.HasPrefix(cmd, "talk "):
		handlers.HandleTalk(conn, cmd, p, npcRepo)
		return false
	case cmd == "buy", strings.HasPrefix(cmd, "buy "):
		handlers.HandleBuy(cmd, conn, p, npcRepo)
		return false
	case cmd == "quest", strings.HasPrefix(cmd, "quest "):
		handlers.HandleQuestList(conn, p)
		return false
	default:
		fmt.Fprintf(conn, "Неизвестная команда\n> ")
		return false
	}

}

// помошник  Разрешённые команды во время диалога с ....
func isCommandBlockedInDialog(cmd string, p *player.Player) bool {
	// Разрешённые команды во время диалога
	if cmd == "stop talk" || cmd == "quit" {
		return false
	}

	// Для торговца
	if cmd == "buy" || strings.HasPrefix(cmd, "buy ") {
		return false
	}

	// Для кузнеца
	if p.TalkToID == "blacksmith" {
		if strings.HasPrefix(cmd, "repair ") || strings.HasPrefix(cmd, "take ") {
			return false
		}
	}

	return true // всё остальное — блокируем
}
