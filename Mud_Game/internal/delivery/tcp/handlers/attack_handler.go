package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// функция атаки на монстра
func HandleAttack(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	//получаем текущую комнату
	room, err := roomRepo.FindByID(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
		return
	}

	monster := room.GetMonster()

	//есть ли в комнате монстр
	if monster == nil {
		fmt.Fprintf(conn, "В комнате монстра нет\n> ")
		return
	}

	//жив ли
	if !monster.IsAlive {
		fmt.Fprintf(conn, "Монстр мертв.\n> ")
		return
	}

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
	finalDamageMonster := damage - monster.Defence
	if finalDamageMonster <= 0 {
		finalDamageMonster = 1
	}

	//наносим урон
	monster.Health -= finalDamageMonster
	room.SetMonster(monster)
	roomRepo.Save(room)

	////////////////////////////////////если монстр умер/////////////////////////////////////
	if monster.Health <= 0 {
		monster.Health = 0
		monster.IsAlive = false
		monster.TimeToLoot = time.Now().Add(40 * time.Second) ///////пока что время на осмотр лута  -
		monster.RespawnTime = time.Now().Add(1 * time.Minute) ///////пока что через 30 сек.
		room.SetMonster(monster)
		roomRepo.Save(room)

		//+опыт
		p.AddExperience(monster.Experience, conn)

		//лут(пока монеты)
		randCoins := 1 + rand.Intn(5)
		player.AddItem(&p.Inventory, "coin", randCoins)

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
				room.ClearItems()
				p.CurrentRoom = "dungeon_entrance_goblins"
				playerRepo.Save(p)
				p.SendMessage(conn, "\n💥 Пещера обвалилась! Тебя выбросило наружу.\n>  ")

				//очищаем occupantID
				room.SetPlayerOccupantID("")
				roomRepo.Save(room)
			}
		}()

		fmt.Fprintf(conn, "Ты нанес %d урона! %s повержен!\n", finalDamageMonster, monster.Name)
		fmt.Fprintf(conn, "Получено %d опыта.\n", monster.Experience)
		fmt.Fprintf(conn, "Найдено %d монет.\n> ", randCoins)
		fmt.Fprintf(conn, "⚠️ Пещера начнёт разрушаться через 2 минуты. У тебя есть время на обыск!\n> ")
		return
	}

	////////////////////////////////////если выжил,он атакует/////////////////////////////////////
	monsterDamage := monster.MinDamage + rand.Intn(monster.MaxDamage-monster.MinDamage+1)

	//учитываем защиту игрока
	defence := p.GetTotalDefence()
	reduction := float64(defence) / (float64(defence) + 100)
	finalDamage := int(float64(monsterDamage) * (1 - reduction))
	if finalDamage <= 0 {
		finalDamage = 1
	}

	p.Stats.Health -= finalDamage

	fmt.Fprintf(conn, "Ты нанес %d урона! %s нанес %d урона!\n", finalDamageMonster, monster.Name, finalDamage)
	fmt.Fprintf(conn, "Здоровье гоблина: %d\n", monster.Health)

	//проверка смерти
	if p.Stats.Health <= 0 {
		fmt.Fprintf(conn, "💀 Ты погиб, персонаж удаляется...\n")
		//очищаем поле "окупанта"
		currentRoom, _ := roomRepo.FindByID(p.CurrentRoom)
		if currentRoom != nil && currentRoom.GetID() == "dungeon_goblin" {
			currentRoom.SetPlayerOccupantID("")
			roomRepo.Save(currentRoom)
		}
		playerRepo.Delete(p.ID)
		conn.Close()
		return
	}
	fmt.Fprintf(conn, "> ")
}
