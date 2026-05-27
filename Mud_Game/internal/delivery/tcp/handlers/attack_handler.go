package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
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
		room.SetMonster(monster)
		roomRepo.Save(room)

		//+опыт
		p.AddExperience(monster.Experience, conn)

		//лут(пока монеты)
		randCoins := 1 + rand.Intn(5)
		player.AddItem(&p.Inventory, "coin", randCoins)

		fmt.Fprintf(conn, "Ты нанес %d урона! %s повержен!\n", finalDamageMonster, monster.Name)
		fmt.Fprintf(conn, "Получено %d опыта.\n", monster.Experience)
		fmt.Fprintf(conn, "Найдено %d монет.\n> ", randCoins)
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
		playerRepo.Delete(p.ID)
		conn.Close()
		return
	}
	fmt.Fprintf(conn, "> ")
}
