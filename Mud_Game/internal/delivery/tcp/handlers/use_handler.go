package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// использование
func HandleUse(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	args, found := strings.CutPrefix(cmd, "use ")
	if !found {
		return
	}
	args = strings.TrimSpace(args)
	if args == "" {
		fmt.Fprintf(conn, "Что использовать? Использование : use <предмет>\n> ")
		return
	}

	parts := strings.Fields(args)
	var itemName string

	if len(parts) == 1 {
		itemName = parts[0]
	} else {
		itemName = strings.Join(parts, " ")
	}

	var index int = -1
	if num, err := strconv.Atoi(itemName); err == nil {
		target, idx := p.FindItemByNumber(num)
		if target == nil {
			fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", num)
			return
		}
		index = idx
		itemName = target.Name
	} else {
		index = p.FindItemIndex(itemName)
	}

	if index == -1 {
		fmt.Fprintf(conn, "Предмет не найден\n> ")
		return
	}
	thatItem := p.Inventory[index]

	if thatItem.ItemType != "scroll" {
		fmt.Fprintf(conn, "Это нельзя использовать\n> ")
		return
	}

	//////////////////////////////////////если свиток "огненый шар"///////////////////////
	if thatItem.Name == "scroll fireball" {
		//получаем текущую комнату
		r, err := roomRepo.FindByID(p.CurrentRoom)
		if err != nil {
			fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
			return
		}

		monster := r.GetMonster()
		//есть ли в комнате монстр
		if monster == nil || !monster.IsAlive {
			fmt.Fprintf(conn, "Нету противников для использования свитка\n> ")
			return
		}

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
			fmt.Fprintf(conn, "💥 Ты прочел заклинание \"Огненного шара\" и свиток превратился в огненный снаряд который ты пустил по противнику и нанес %d урона.\n> ", finalDamage)
			if monster.Health <= 0 {
				monster.Health = 0
				monster.IsAlive = false
				fmt.Fprintf(conn, "Монстр повержен!\n")
			}
			r.SetMonster(monster)
			roomRepo.Save(r)
		} else {
			fmt.Fprintf(conn, "Ты попытался использовать заклинание, но случайно сжег свиток.\n> ")
		}
		player.RemoveItem(&p.Inventory, thatItem.Name, 1)
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
		player.RemoveItem(&p.Inventory, thatItem.Name, 1)
		playerRepo.Save(p)
		return

	}

}
