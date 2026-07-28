package handlers

import (
	"Mud_game/Mud_Game/internal/domain/monster"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// функция атаки на монстра
func HandleAttack(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	//получаем текущую комнату
	r, err := roomRepo.FindByID(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
		return
	}

	//приведение комнаты в интерфейс
	concreteRoom, ok := r.(*room.Room)
	if !ok {
		fmt.Fprintf(conn, "Ошибка приведения комнаты\n> ")
		return
	}

	// Получаем список живых монстров
	allAliveMonster := concreteRoom.GetAliveMonsters()

	if len(allAliveMonster) == 0 {

		fmt.Fprintf(conn, "В комнате нет монстров\n> ")
		return
	}

	// Разбираем команду
	args := strings.Fields(cmd)
	var targetIndex int = 0

	if len(args) > 1 {
		idx, err := strconv.Atoi(args[1])
		if err == nil && idx >= 1 && idx <= len(allAliveMonster) {
			targetIndex = idx - 1
		} else {
			fmt.Fprintf(conn, "Некорректный номер цели. Доступно: 1-%d\n> ", len(allAliveMonster))
			return
		}
	}

	// ✅ РАБОТАЕМ НАПРЯМУЮ С allAliveMonster[targetIndex]
	selectedMonster := allAliveMonster[targetIndex]

	p.Stats.IsInCombat = true

	//Проверка яда ПОСЛЕ выбора цели
	if p.Stats.IsPoisoned {
		p.Stats.Health -= selectedMonster.PoisonDamage
		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб от яда..\n")
			playerRepo.Delete(p.ID)
			conn.Close()
			return
		}
	}

	// ... наносим урон monster ...

	//рассчитываем урон игрока
	weapon := p.Equipment.Weapon
	minDamage := 1
	maxDamage := 3

	if weapon != nil {
		minDamage = weapon.MinDamage
		maxDamage = weapon.MaxDamage
	}

	//бонус к урону от силы +1 за каждые 3 очка силы
	strengthBonus := p.Stats.Strength / 3

	damage := minDamage + rand.Intn(maxDamage-minDamage+1) + strengthBonus //не учли броню?

	//учитываем броню монстра
	finalDamageMonster := damage - selectedMonster.Defence
	if finalDamageMonster <= 0 {
		finalDamageMonster = 1
	}

	var targetMonster *monster.Monster

	// Сначала ищем в MonsterS (новый данж)
	for i, m := range concreteRoom.MonsterS {
		if m.ID == selectedMonster.ID {

			targetMonster = concreteRoom.MonsterS[i]
			break
		}
	}

	// Если не нашли — ищем в Monster (старый данж)
	if targetMonster == nil && concreteRoom.Monster != nil {
		if concreteRoom.Monster.ID == selectedMonster.ID {
			targetMonster = concreteRoom.Monster

		}
	}

	if targetMonster == nil {

		fmt.Fprintf(conn, "Ошибка: монстр не найден\n> ")
		return
	}

	// Наносим урон
	targetMonster.Health -= finalDamageMonster
	selectedMonster = targetMonster

	if targetMonster.Health <= 0 {
		targetMonster.IsAlive = false
		selectedMonster.IsAlive = false
	}
	// Обновляем монстра в комнате
	if len(concreteRoom.MonsterS) == 0 {
		concreteRoom.Monster = selectedMonster
	} else {
		// Для нового данжа обновляем MonsterS
		for i, m := range concreteRoom.MonsterS {
			if m.ID == selectedMonster.ID {
				concreteRoom.MonsterS[i].Health = selectedMonster.Health
				concreteRoom.MonsterS[i].IsAlive = selectedMonster.IsAlive
				break
			}
		}
	}
	roomRepo.Save(concreteRoom)

	////////////////////////////////////если монстр умер/////////////////////////////////////
	if selectedMonster.Health <= 0 {
		selectedMonster.Health = 0
		selectedMonster.IsAlive = false

		// ✅ Проверяем, все ли монстры мертвы
		allDead := true
		for _, m := range concreteRoom.MonsterS {
			if m.IsAlive {
				allDead = false
				break
			}
		}

		// Если есть ещё живые — обвал НЕ наступает
		if !allDead {

			fmt.Fprintf(conn, "====================\n> ")

		} else {

			// После того как все монстры убиты - обвал
			if p.Stats.IsPoisoned {
				go p.StartPoisonTicker(conn, playerRepo)
			}

			p.Stats.IsInCombat = false

			// ✅ Только если все мертвы — запускаем обвал
			p.StopDungeonTimer()
			selectedMonster.TimeToLoot = time.Now().Add(40 * time.Second)
			selectedMonster.RespawnTime = time.Now().Add(1 * time.Minute)

			// Добавляем выход
			if concreteRoom.Exits == nil {
				concreteRoom.Exits = make(map[string]string)
			}
			concreteRoom.Exits["up"] = "dungeon_entrance_goblins"

			concreteRoom.SetMonster(selectedMonster)
			roomRepo.Save(concreteRoom)

			//+опыт
			p.AddExperience(selectedMonster.Experience, conn)

			//предупреждение об обвале за 20 секунд
			go func() {
				time.Sleep(20 * time.Second)
				if p.CurrentRoom == "dungeon_goblin" {
					p.SendMessage(conn, "\n💥 Стены пещеры сильно трясутся! Камни падают с потолка! Пещера вот-вот обвалится!\n> ")
				}
			}()

			//время на осмотр лута
			go func() {
				time.Sleep(40 * time.Second) //////пока что время на осмотр лута  (обратный отсчет)
				//проверяем что игрок еще в данже
				if p.CurrentRoom == "dungeon_goblin" {
					//телепортируем на вход и уничтожаем если не успели забрать предметы
					concreteRoom.ClearItems()
					p.CurrentRoom = "dungeon_entrance_goblins"
					playerRepo.Save(p)
					p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n>  ")

					//очищаем occupantID
					concreteRoom.SetPlayerOccupantID("")
					roomRepo.Save(concreteRoom)
				}
			}()

			fmt.Fprintf(conn, "Ты нанес %d урона %s\n", finalDamageMonster, selectedMonster.Name)
			fmt.Fprintf(conn, "Получено %d опыта.\n", selectedMonster.Experience)
			fmt.Fprintf(conn, "⚠️ Пещера начнёт разрушаться через 2 минуты. У тебя есть время на обыск!\n> ")
			return
		}
	}

	////////////////////////////////////они атакуют/////////////////////////////////////

	var actions []string
	var statuses []string

	//ОТВЕТНАЯ АТАКА ВСЕХ МОНСТРОВ//
	for i, m := range allAliveMonster {
		if m.IsAlive {
			statuses = append(statuses, fmt.Sprintf("%d. %s %d/%d HP", i+1, m.Name, m.Health, m.MaxHealth))
		} else {
			statuses = append(statuses, fmt.Sprintf("%d. %s мертв💀", i+1, m.Name))
		}

		if !m.IsAlive {
			continue
		}

		///////////////если шаман гоблин://///
		if m.ID == "goblin_shaman" {
			if m.CastTime < 2 {
				//колдует - яда нет
				m.CastTime++
				//обновляем CastTime в оригинале монстра
				for j, orig := range concreteRoom.MonsterS {
					if orig.ID == m.ID {
						concreteRoom.MonsterS[j].CastTime = m.CastTime
						break
					}
				}

				actions = append(actions, fmt.Sprintf("%d. %s колдует! (%d/2)", i+1, m.Name, m.CastTime))

				continue
			} else {
				//после двух ходов бьет посохом и накладывает яд
				//+яд:
				m.IsCasting = false

				if m.CastTime == 2 {
					p.Stats.IsPoisoned = true
					p.Stats.PoisonTicks = 6               //6 тиков
					p.Stats.PoisonDamage = m.PoisonDamage //берем из данного монстра
					m.CastTime = 3

					//обновляем CastTime в оригинале монстра
					for j, orig := range concreteRoom.MonsterS {
						if orig.ID == m.ID {
							concreteRoom.MonsterS[j].CastTime = m.CastTime
							break
						}
					}
					actions = append(actions, fmt.Sprintf("%d. %s отравил тебя! Яд будет наносить %d урона каждый ход.", i+1, m.Name, m.PoisonDamage))

				}

				//+урон посохом:

				monsterDamage := m.MinDamage + rand.Intn(m.MaxDamage-m.MinDamage+1)

				//учитываем защиту игрока
				defence := p.GetTotalDefence()
				reduction := float64(defence) / (float64(defence) + 100)
				finalDamage := int(float64(monsterDamage) * (1 - reduction))
				if finalDamage <= 0 {
					finalDamage = 1
				}

				p.Stats.Health -= finalDamage

				actions = append(actions, fmt.Sprintf("%d. %s нанес %d урона", i+1, m.Name, finalDamage))
			}

		} else {

			///////////////если физик какой то://////////////////
			monsterDamage := m.MinDamage + rand.Intn(m.MaxDamage-m.MinDamage+1)

			//учитываем защиту игрока
			defence := p.GetTotalDefence()
			reduction := float64(defence) / (float64(defence) + 100)
			finalDamage := int(float64(monsterDamage) * (1 - reduction))
			if finalDamage <= 0 {
				finalDamage = 1
			}

			p.Stats.Health -= finalDamage

			actions = append(actions, fmt.Sprintf("%d. %s нанес %d урона", i+1, m.Name, finalDamage))
		}
	}
	// Выводим результат
	fmt.Fprintf(conn, "Ты нанес %d урона %s\n", finalDamageMonster, selectedMonster.Name)
	for _, status := range statuses {
		fmt.Fprintf(conn, "%s\n", status)
	}
	fmt.Fprintf(conn, "================================\n")
	for _, action := range actions {
		fmt.Fprintf(conn, "%s\n", action)
	}

	if p.Stats.IsPoisoned {
		fmt.Fprintf(conn, "💀 Яд нанес тебе %d урона!\n", p.Stats.PoisonDamage)
	}

	roomRepo.Save(concreteRoom)

	//проверка смерти
	if p.Stats.Health <= 0 {
		//шанс выжить:
		if p.Stats.Tracking >= 6 {
			chance := p.Stats.Tracking * 20 //////////////////////для теста
			if rand.Intn(100) < chance {
				//если выжил
				p.Stats.Health = 1
				p.CurrentRoom = concreteRoom.GetExitRoomID()
				p.BreakAllEquipment()
				concreteRoom.SetPlayerOccupantID("")
				// ✅ ВОССТАНАВЛИВАЕМ ВСЕХ МОНСТРОВ
				if len(concreteRoom.MonsterS) > 0 {
					for _, m := range concreteRoom.MonsterS {
						m.Health = m.MaxHealth
						m.IsAlive = true
						if m.ID == "goblin_shaman" {
							m.CastTime = 0
							m.IsCasting = true
						}
					}
				} else if concreteRoom.Monster != nil {
					concreteRoom.Monster.Health = concreteRoom.Monster.MaxHealth
					concreteRoom.Monster.IsAlive = true
					concreteRoom.Monster.CastTime = 0
				}

				roomRepo.Save(concreteRoom)

				// Яд все равно остается
				if p.Stats.IsPoisoned {
					go p.StartPoisonTicker(conn, playerRepo)
				}

				p.Stats.IsInCombat = false

				playerRepo.Save(p)
				fmt.Fprintf(conn, "🔥 Ты чудом выжил! Инстинкты спасли тебя.\n")
				fmt.Fprintf(conn, "Твоё снаряжение повреждено!\n")
				fmt.Fprintf(conn, "Ты оказался на входе в подземелье.\n> ")

				//шанс потерять чтото из инвентаря при побеге
				if len(p.Inventory) > 0 {
					//⚠️ Важно: Идём с конца, чтобы не сбивать индексы при удалении.
					for i := len(p.Inventory) - 1; i >= 0; i-- {
						if rand.Intn(100) < 30 {
							lostItem := p.Inventory[i]
							p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
							fmt.Fprintf(conn, "При побеге ты потерял %s\n", lostItem.Name)
						}
					}
				}
				return

			}
		}
		//восстановление подземелья/удаление персонажа
		// ✅ ВОССТАНАВЛИВАЕМ ВСЕХ МОНСТРОВ
		if len(concreteRoom.MonsterS) > 0 {
			for _, m := range concreteRoom.MonsterS {
				m.Health = m.MaxHealth
				m.IsAlive = true
				if m.ID == "goblin_shaman" {
					m.CastTime = 0
					m.IsCasting = true
				}
			}
		} else if concreteRoom.Monster != nil {
			concreteRoom.Monster.Health = concreteRoom.Monster.MaxHealth
			concreteRoom.Monster.IsAlive = true
			concreteRoom.Monster.CastTime = 0
		}

		concreteRoom.SetPlayerOccupantID("")
		roomRepo.Save(concreteRoom)
		fmt.Fprintf(conn, "💀 Ты погиб, персонаж удаляется...\n")
		playerRepo.Delete(p.ID) //////////////////
		conn.Close()
		return
	}

	fmt.Fprintf(conn, "> ")

}
