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

	// 1. Получение комнаты и монстров
	concreteRoom, allAliveMonster := getRoomAndMonsters(conn, p, roomRepo)
	if concreteRoom == nil {
		return
	}

	// 2. Выбор цели
	selectedMonster, _ := selectTarget(conn, cmd, allAliveMonster)
	if selectedMonster == nil {
		return
	}

	// 3. Проверка яда
	if applyPoisonDamage(conn, p, selectedMonster, playerRepo) {
		return
	}

	// 4. Урон игрока
	damage, fireDamage, magicDamage, poisonDamage, finalDamage := calculatePlayerDamage(p, selectedMonster)
	// 5. Нанесение урона
	targetMonster := allpyDamageToMonster(concreteRoom, selectedMonster, finalDamage)
	if targetMonster == nil {
		fmt.Fprintf(conn, "Ошибка: монстр не найден\n> ")
		return
	}
	selectedMonster = targetMonster
	// 6. Проверка смерти монстра
	if selectedMonster.Health <= 0 {
		if handleMonsterDeath(conn, selectedMonster, concreteRoom, playerRepo, p, damage) {
			roomRepo.Save(concreteRoom)
			return //обвал
		}
	}
	// 7. Ответная атака монстра если жив
	actions, statuses := otvetAttakaMonstra(allAliveMonster, concreteRoom, p)

	// 8. Вывод результата
	printBattleResult(conn, damage, selectedMonster, statuses, actions, p, fireDamage, magicDamage, poisonDamage)
	roomRepo.Save(concreteRoom)
	// 9. Проверка смерти игрока
	proverkaSmerti(conn, concreteRoom, p, roomRepo, playerRepo)
	fmt.Fprintf(conn, "> ")

}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// получение комнаты и монстров
func getRoomAndMonsters(conn net.Conn, p *player.Player, roomRepo room.Repository) (*room.Room, []*monster.Monster) {
	//получаем текущую комнату
	r, err := roomRepo.FindByID(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
		return nil, nil
	}

	//приведение комнаты в интерфейс
	concreteRoom, ok := r.(*room.Room)
	if !ok {
		fmt.Fprintf(conn, "Ошибка приведения комнаты\n> ")
		return nil, nil
	}

	// Получаем список живых монстров
	allAliveMonster := concreteRoom.GetAliveMonsters()

	if len(allAliveMonster) == 0 {

		fmt.Fprintf(conn, "В комнате нет монстров\n> ")
		return nil, nil
	}
	return concreteRoom, allAliveMonster
}

// выбор цели
func selectTarget(conn net.Conn, cmd string, allAliveMonster []*monster.Monster) (*monster.Monster, int) {
	// Разбираем команду
	args := strings.Fields(cmd)
	var targetIndex int = 0

	if len(args) > 1 {
		idx, err := strconv.Atoi(args[1])
		if err == nil && idx >= 1 && idx <= len(allAliveMonster) {
			targetIndex = idx - 1
		} else {
			fmt.Fprintf(conn, "Некорректный номер цели. Доступно: 1-%d\n> ", len(allAliveMonster))
			return nil, -1
		}
	}
	return allAliveMonster[targetIndex], targetIndex
}

// проверка яда
func applyPoisonDamage(conn net.Conn, p *player.Player, selectedMonster *monster.Monster, playerRepo player.Repository) bool {

	if p.Stats.IsPoisoned {
		p.Stats.Health -= selectedMonster.PoisonDamage
		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб от яда..\n")
			playerRepo.Delete(p.ID)
			conn.Close()
			return true
		}
	}
	return false
}

// расчитываем  урон игрока
func calculatePlayerDamage(p *player.Player, monster *monster.Monster) (int, int, int, int, int) {

	weapon := p.Equipment.Weapon
	minDamage := 1
	maxDamage := 3

	if weapon != nil {
		minDamage = weapon.MinDamage
		maxDamage = weapon.MaxDamage
	}

	//бонус к урону от силы +1 за каждые 3 очка силы
	strengthBonus := p.Stats.Strength / 3
	baseDamage := minDamage + rand.Intn(maxDamage-minDamage+1) + strengthBonus //БАЗОВЫЙ урон оружия(или руки)

	// 🔥 Огненный урон
	fireDamage := 0
	if weapon != nil && weapon.FireDamage > 0 {
		fireDamage = weapon.FireDamage
		reduction := float64(monster.FireDefence) / (float64(monster.FireDefence) + 100)
		fireDamage = int(float64(fireDamage) * (1 - reduction))
	}
	// 🧙 Магический урон
	magicDamage := 0
	if weapon != nil && weapon.MagicDamage > 0 {
		magicDamage = weapon.MagicDamage
		reduction := float64(monster.MagicDefence) / (float64(monster.MagicDefence) + 100)
		magicDamage = int(float64(magicDamage) * (1 - reduction))
	}
	// ☠️ Ядовитый урон
	poisonDamage := 0
	if weapon != nil && weapon.PoisonDamage > 0 {
		poisonDamage = weapon.PoisonDamage
		reduction := float64(monster.PoisonDefence) / (float64(monster.PoisonDefence) + 100)
		poisonDamage = int(float64(poisonDamage) * (1 - reduction))
	}
	//физ урон
	reduction := float64(monster.Defence) / (float64(monster.Defence) + 100)
	physicalDamage := int(float64(baseDamage) * (1 - reduction))

	finalDamage := physicalDamage + fireDamage + poisonDamage + magicDamage
	if finalDamage <= 0 {
		finalDamage = 1
	}

	return physicalDamage, fireDamage, magicDamage, poisonDamage, finalDamage
}

// наносим урон монстру
func allpyDamageToMonster(concreteRoom *room.Room, selectedMonster *monster.Monster, damage int) *monster.Monster {
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
		return nil
	}

	// Наносим урон
	targetMonster.Health -= damage
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

	return targetMonster
}

// проверка смерти монстра
func handleMonsterDeath(conn net.Conn, selectedMonster *monster.Monster, concreteRoom *room.Room, playerRepo player.Repository, p *player.Player, damage int) bool {
	if selectedMonster.Health <= 0 {
		selectedMonster.Health = 0
		selectedMonster.IsAlive = false

		//+опыт
		p.AddExperience(selectedMonster.Experience, conn)

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
			return false

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

			if p.CurrentRoom == "dungeon_goblin" {
				concreteRoom.Exits["up"] = "dungeon_entrance_goblins"
			}

			if p.CurrentRoom == "dungeon_goblins_v2" {
				concreteRoom.Exits["up"] = "dungeon_entrance_goblins_v2"
				concreteRoom.Exits["down"] = "glubini_room"
			}

			if p.CurrentRoom == "glubini_room" {
				concreteRoom.Exits["up"] = "dungeon_entrance_goblins_v2"
			}
			concreteRoom.SetMonster(selectedMonster)

			//предупреждение об обвале за 20 секунд
			go func() {
				time.Sleep(20 * time.Second)
				if p.CurrentRoom == "dungeon_goblin" ||
					p.CurrentRoom == "dungeon_goblins_v2" ||
					p.CurrentRoom == "glubini_room" {
					p.SendMessage(conn, "\n💥 Стены пещеры сильно трясутся! Камни падают с потолка! Пещера вот-вот обвалится!\n> ")
				}
			}()

			//время на осмотр лута
			go func() {
				time.Sleep(40 * time.Second) //////пока что время на осмотр лута  (обратный отсчет)
				//проверяем что игрок еще в данже
				//У ГОБЛИНА
				if p.CurrentRoom == "dungeon_goblin" {
					//телепортируем на вход и уничтожаем если не успели забрать предметы
					concreteRoom.ClearItems()
					p.CurrentRoom = "dungeon_entrance_goblins"
					playerRepo.Save(p)
					p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n>  ")

					//очищаем occupantID
					concreteRoom.SetPlayerOccupantID("")

				}

				//У ДВУХ ГОБЛИНОВ
				if p.CurrentRoom == "dungeon_goblins_v2" {
					//телепортируем на вход и уничтожаем если не успели забрать предметы
					concreteRoom.ClearItems()
					p.CurrentRoom = "dungeon_entrance_goblins_v2"
					playerRepo.Save(p)
					p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n>  ")

					//очищаем occupantID
					concreteRoom.SetPlayerOccupantID("")

				}

				//В ГЛУБИНАХ
				if p.CurrentRoom == "glubini_room" {
					//телепортируем на вход и уничтожаем если не успели забрать предметы
					concreteRoom.ClearItems()
					p.CurrentRoom = "dungeon_entrance_goblins_v2"
					playerRepo.Save(p)
					p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n>  ")

					//очищаем occupantID
					concreteRoom.SetPlayerOccupantID("")

				}
			}()

			fmt.Fprintf(conn, "Ты нанес %d урона %s\n", damage, selectedMonster.Name)
			fmt.Fprintf(conn, "Получено %d опыта.\n", selectedMonster.Experience)
			fmt.Fprintf(conn, "⚠️ Пещера начнёт разрушаться через 2 минуты. У тебя есть время на обыск!\n> ")
			return true
		}
	}
	return true
}

// ответка от монстра если выжил
func otvetAttakaMonstra(allAliveMonster []*monster.Monster, concreteRoom *room.Room, p *player.Player) ([]string, []string) {
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
			/////////////////если  ВЕРХОВНЫЙ шаман гоблин://///////////////////////////////
		} else if m.ID == "goblin_high_shaman" {
			monsterDamage := m.MagicDamage + rand.Intn(4)
			defence := p.GetTotalMagicDefence()
			reduction := float64(defence) / (float64(defence) + 100)
			finalDamage := int(float64(monsterDamage) * (1 - reduction))
			if finalDamage <= 0 {
				finalDamage = 1
			}
			p.Stats.Health -= finalDamage

			actions = append(actions, fmt.Sprintf("%d. %s нанес %d магического урона", i+1, m.Name, finalDamage))
			continue
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
	return actions, statuses

}

// вывод результата хода
func printBattleResult(conn net.Conn, damage int, selectedMonster *monster.Monster, statuses []string, actions []string, p *player.Player, fireDamage int, magicDamage int, poisonDamage int) {
	// Основная строка урона
	fmt.Fprintf(conn, "Ты нанес %d урона %s\n", damage, selectedMonster.Name)

	// Бонусы на отдельной строке с отступом
	var bonuses []string
	if fireDamage > 0 {
		bonuses = append(bonuses, fmt.Sprintf("+%d огненного урона", fireDamage))
	}
	if magicDamage > 0 {
		bonuses = append(bonuses, fmt.Sprintf("+%d магического урона", magicDamage))
	}
	if poisonDamage > 0 {
		bonuses = append(bonuses, fmt.Sprintf("+%d ядовитого урона", poisonDamage))
	}

	if len(bonuses) > 0 {
		fmt.Fprintf(conn, "           (%s)\n", strings.Join(bonuses, ", "))
	}

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

}

// проверка смерти
func proverkaSmerti(conn net.Conn, concreteRoom *room.Room, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
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
}
