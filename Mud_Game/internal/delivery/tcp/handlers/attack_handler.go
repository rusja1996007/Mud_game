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

	// Находим реальный индекс монстра в оригинальном списке
	var monsterIndex int = -1
	var targetMonster *monster.Monster

	// Сначала ищем в MonsterS (новый данж)
	for i, m := range concreteRoom.MonsterS {
		if m.ID == allAliveMonster[targetIndex].ID {
			monsterIndex = i
			targetMonster = concreteRoom.MonsterS[i]
			break
		}
	}

	// Если не нашли — ищем в Monster (старый данж)
	if targetMonster == nil && concreteRoom.Monster != nil {
		if concreteRoom.Monster.ID == allAliveMonster[targetIndex].ID {
			targetMonster = concreteRoom.Monster
			monsterIndex = 0 // фиктивный индекс, не используется для MonsterS
		}
	}
	// monsterIndex не используется дальше, но он нужен для компиляции
	// Можешь добавить в конце:
	_ = monsterIndex // ← если компилятор ругается

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
			fmt.Fprintf(conn, "Ты нанес %d урона! %s повержен!\n> ", finalDamageMonster, selectedMonster.Name)
			roomRepo.Save(concreteRoom)
			return
		}
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

		fmt.Fprintf(conn, "Ты нанес %d урона! %s повержен!\n", finalDamageMonster, selectedMonster.Name)
		fmt.Fprintf(conn, "Получено %d опыта.\n", selectedMonster.Experience)
		fmt.Fprintf(conn, "⚠️ Пещера начнёт разрушаться через 2 минуты. У тебя есть время на обыск!\n> ")
		return
	}

	////////////////////////////////////если выжил,он атакует/////////////////////////////////////

	///////////////если шаман гоблин://///
	if selectedMonster.ID == "goblin_shaman" {
		if selectedMonster.CastTime < 2 {
			selectedMonster.CastTime++

			fmt.Fprintf(conn, "🧙Гоблин шаман колдует! (%d/2)\n", selectedMonster.CastTime)
			fmt.Fprintf(conn, "Ты нанес %d урона.\n", finalDamageMonster)
			fmt.Fprintf(conn, "Здоровье гоблина шамана: %d\n", selectedMonster.Health)

			//обновляем CastTime в оригинале монстра
			for i, m := range concreteRoom.MonsterS {
				if m.ID == selectedMonster.ID {
					concreteRoom.MonsterS[i].CastTime = selectedMonster.CastTime
					break
				}
			}
			roomRepo.Save(concreteRoom)
			fmt.Fprintf(conn, "> ")
			return
		} else {
			//после двух ходов бьет посохом
			selectedMonster.IsCasting = false
			monsterDamage := selectedMonster.MinDamage + rand.Intn(selectedMonster.MaxDamage-selectedMonster.MinDamage+1)

			//учитываем защиту игрока
			defence := p.GetTotalDefence()
			reduction := float64(defence) / (float64(defence) + 100)
			finalDamage := int(float64(monsterDamage) * (1 - reduction))
			if finalDamage <= 0 {
				finalDamage = 1
			}

			p.Stats.Health -= finalDamage

			fmt.Fprintf(conn, "Ты нанес %d урона! %s нанес %d урона!\n", finalDamageMonster, selectedMonster.Name, finalDamage)
			fmt.Fprintf(conn, "Здоровье гоблина шамана: %d\n", selectedMonster.Health)

		}

	} else {

		///////////////если физик какой то://////////////////
		monsterDamage := selectedMonster.MinDamage + rand.Intn(selectedMonster.MaxDamage-selectedMonster.MinDamage+1)

		//учитываем защиту игрока
		defence := p.GetTotalDefence()
		reduction := float64(defence) / (float64(defence) + 100)
		finalDamage := int(float64(monsterDamage) * (1 - reduction))
		if finalDamage <= 0 {
			finalDamage = 1
		}

		p.Stats.Health -= finalDamage

		fmt.Fprintf(conn, "Ты нанес %d урона! %s нанес %d урона!\n", finalDamageMonster, selectedMonster.Name, finalDamage)
		fmt.Fprintf(conn, "Здоровье гоблина: %d\n", selectedMonster.Health)
	}

	//проверка смерти
	if p.Stats.Health <= 0 {
		//шанс выжить:
		if p.Stats.Tracking >= 6 {
			chance := p.Stats.Tracking * 14 //////////////////////для теста
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
				}

				roomRepo.Save(concreteRoom)
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
		}
		concreteRoom.SetPlayerOccupantID("")

		roomRepo.Save(concreteRoom)
		fmt.Fprintf(conn, "💀 Ты погиб, персонаж удаляется...\n")
		p.StopAllTickers()
		playerRepo.Delete(p.ID) //////////////////
		conn.Close()
		return
	}
	fmt.Fprintf(conn, "> ")
}
