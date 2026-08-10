package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// использование
func HandleUse(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	args, found := strings.CutPrefix(cmd, "use ")
	if !found {
		return
	}
	parts := strings.Fields(args)

	var itemName string
	var index int = -1
	var inBag bool

	// Проверяем, что есть аргументы
	if len(parts) == 0 {
		fmt.Fprintf(conn, "Что использовать? Использование: use <номер_свитка> <цель>\n> ")
		return
	}

	// Определяем номер предмета (всегда первый аргумент)
	itemNum, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Fprintf(conn, "Нужно указать номер предмета\n> ")
		return
	}

	// Ищем предмет по номеру
	target, idx := p.FindItemByNumber(itemNum)
	if target == nil {
		fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", itemNum)
		return
	}
	itemName = target.Name
	index = idx
	if itemNum <= len(p.Inventory) {
		inBag = false
	} else {
		inBag = true
	}

	// Получаем предмет
	var thatItem *item.ItemStack
	if inBag {
		thatItem = p.Equipment.BagItems[index]
	} else {
		thatItem = p.Inventory[index]
	}

	// Проверяем тип
	if thatItem.ItemType != "scroll" {
		fmt.Fprintf(conn, "Это нельзя использовать\n> ")
		return
	}

	//////////////////////////////////////если свиток "огненый шар"///////////////////////
	if thatItem.Name == "scroll fireball" {
		// Проверяем, что указана цель
		if len(parts) < 2 {
			fmt.Fprintf(conn, "Укажите цель: use <свиток> <цель>\n> ")
			return
		}

		targetMonsterNum, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Fprintf(conn, "Некорректный номер цели\n> ")
			return
		}

		//получаем текущую комнату
		r, err := roomRepo.FindByID(p.CurrentRoom)
		if err != nil {
			fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
			return
		}

		concreteRoom, ok := r.(*room.Room)
		if !ok {
			fmt.Fprintf(conn, "Ошибка приведения комнаты\n> ")
			return
		}
		// Получаем живых монстров
		allAlive := concreteRoom.GetAliveMonsters()
		if len(allAlive) == 0 {
			fmt.Fprintf(conn, "Нету противников для использования свитка\n> ")
			return
		}

		// Проверяем номер цели
		if targetMonsterNum < 1 || targetMonsterNum > len(allAlive) {
			fmt.Fprintf(conn, "Некорректный номер цели. Доступно: 1-%d\n> ", len(allAlive))
			return
		}

		monster := allAlive[targetMonsterNum-1]

		damage := thatItem.MinDamage + rand.Intn(thatItem.MaxDamage-thatItem.MinDamage+1)
		finalDamage := damage - monster.FireDefence
		if finalDamage < 1 {
			finalDamage = 1
		}

		uspeh := true

		//если интелект меньше либо равно 4
		if p.Stats.Intelect <= 4 {
			if rand.Intn(100) >= 50 {
				uspeh = false
			}
		}

		if uspeh {
			monster.Health -= finalDamage
			fmt.Fprintf(conn, "💥 Ты прочел заклинание \"Огненного шара\" и свиток превратился в огненный снаряд который ты пустил по противнику и нанес %d урона.\n", finalDamage)
			if monster.Health <= 0 {
				monster.Health = 0
				monster.IsAlive = false
				fmt.Fprintf(conn, "🔥 %s повержен!\n", monster.Name)

				allDead := true
				for _, m := range concreteRoom.MonsterS {
					if m.IsAlive {
						allDead = false
						break
					}
				}
				if allDead {
					/////////ЛОГИКА ОБВАЛА КАК В handleMonsterDeath/////////////////
					// ✅ ВСЕ МОНСТРЫ МЕРТВЫ → ОБВАЛ
					if p.Stats.IsPoisoned {
						go p.StartPoisonTicker(conn, playerRepo)
					}

					p.Stats.IsInCombat = false
					p.StopDungeonTimer()

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
						concreteRoom.Exits["up"] = "dungeon_goblins_v2"
					}

					concreteRoom.SetMonster(monster)

					//тАЙМЕР ОБВАЛА
					go func() {
						time.Sleep(40 * time.Second)
						//проверяем что игрок еще в этой комнате
						if p.CurrentRoom == "dungeon_goblin" ||
							p.CurrentRoom == "dungeon_goblins_v2" ||
							p.CurrentRoom == "glubini_room" {
							concreteRoom.ClearItems()
							var exiRoom string
							switch p.CurrentRoom {
							case "dungeon_goblin":
								exiRoom = "dungeon_entrance_goblins"
							case "dungeon_goblins_v2", "glubini_room":
								exiRoom = "dungeon_entrance_goblins_v2"
							}
							p.CurrentRoom = exiRoom
							playerRepo.Save(p)
							p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n> ")
							concreteRoom.SetPlayerOccupantID("")
							roomRepo.Save(concreteRoom)
						}
					}()
					fmt.Fprintf(conn, "⚠️ Пещера начнёт разрушаться через 2 минуты. У тебя есть время на обыск!\n> ")
				} else {
					// ✅ Если есть живые монстры
					fmt.Fprintf(conn, "> ")
				}
			} else {
				fmt.Fprintf(conn, "> ")
			}
			r.SetMonster(monster)
			roomRepo.Save(r)
		} else {
			fmt.Fprintf(conn, "Ты попытался использовать заклинание, но случайно сжег свиток.\n> ")
		}
		p.RemoveOneItem(itemName, inBag, index)
		playerRepo.Save(p)
		return
	}

	/////////////////////////////////////////////////////если свиток "хил"///////////////////////
	if thatItem.Name == "scroll heal" {

		heal := thatItem.HealMin + rand.Intn(thatItem.HealMax-thatItem.HealMin+1)

		maxHealt := 50 + p.Stats.Strength*5 + p.Stats.MaxHealthBonus
		uspeh := true
		if p.Stats.Intelect <= 4 {
			if rand.Intn(100) >= 50 {
				uspeh = false
			}
		}

		if uspeh {
			p.Stats.Health += heal
			if p.Stats.Health >= maxHealt {
				p.Stats.Health = maxHealt
			}
			fmt.Fprintf(conn, "Ты прочел заклинание \"Исцеления\" он превратился в свет и восстановил тебе %d здоровья\n> ", heal)

		} else {
			fmt.Fprintf(conn, "Ты попытался использовать заклинание, но случайно сжег свиток.\n> ")
		}
		p.RemoveOneItem(itemName, inBag, index)
		playerRepo.Save(p)
		return

	}

}
