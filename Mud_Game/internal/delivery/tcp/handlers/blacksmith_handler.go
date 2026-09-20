package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/npc"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type BrokenRef struct {
	Index int  //индекс в p.Inventory или мешке(p.Equipment.BagItems)
	InBag bool //truy - мешок, false - инвентарь
}

// взаимодействие с кузнецом
func HandleBlacksmith(conn net.Conn, p *player.Player, target *npc.NPC, npcRepo *npc_repo.PostgresNPCRepository) {
	p.IsTalkin = true
	p.TalkToID = target.ID
	p.TalkToName = target.Name

	//Собрать список поврежденных предметов сначало из инвентаря
	fmt.Fprintf(conn, "🔨 Кузнец: \nДобрый день! Что будем чинить?\n")
	var brokenItems []BrokenRef // индексы в p.Inventory
	canRepair := map[string]bool{
		"weapon": true,
		"armor":  true,
		"helmet": true,
		"shield": true,
		"boots":  true,
		"bag":    true,
		"ring":   true,
	}
	for i, item := range p.Inventory {
		if item.Durability < 100 && canRepair[item.ItemType] {
			brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: false})
		}
	}

	//.. теперь из мешка если есть
	if p.Equipment.Bag != nil {
		for i, item := range p.Equipment.BagItems {
			if item.Durability < 100 && canRepair[item.ItemType] {
				brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: true})
			}
		}
	}

	//если пустой списко
	if len(brokenItems) == 0 {
		fmt.Fprintf(conn, "Тебе нечего чинить\n> ")
		return
	}

	counter := 0 // общий счетчик

	//показать сам список
	fmt.Fprintf(conn, "📋Твои предметы:\n")
	for _, ref := range brokenItems {
		counter++
		var items *item.ItemStack
		if ref.InBag {
			items = p.Equipment.BagItems[ref.Index]
		} else {
			items = p.Inventory[ref.Index]
		}
		price := npc.CalculateRepairPrice(items)   //цена
		duration := npc.CalculateRepairTime(items) //время
		fmt.Fprintf(conn, "%d. %s(Прочность %d/100) - ремонт %d coins, %v\n", counter, item.GetColoredName(items), items.Durability, price, duration)
	}

	//показать которые в ремонте(только свои)
	// Сначала посчитать свои заказы
	myOrdersCount := 0
	for _, order := range target.RepairOrders {
		if order.PlayerID == p.ID {
			myOrdersCount++
		}
	}

	if myOrdersCount > 0 {
		fmt.Fprintf(conn, "🔨 В ремонте:\n")
		for _, order := range target.RepairOrders {
			if order.PlayerID == p.ID {
				counter++
				if time.Now().After(order.EndTime) {
					fmt.Fprintf(conn, "%d. %s - ГОТОВ(take %d)\n", counter, item.GetColoredName(order.Item), counter)
				} else {
					fmt.Fprintf(conn, "%d. %s - будет готов через %v\n", counter, order.Item.Name, time.Until(order.EndTime).Round(time.Second))
				}

			}

		}
	}
	fmt.Fprintf(conn, "Используй 'repair <номер>' чтобы починить\n")
	fmt.Fprintf(conn, "Используй 'take <номер>' чтобы забрать готовый предмет\n")
	fmt.Fprintf(conn, "Используй 'stop talk' чтобы закончить разговор\n> ")
}

// починка (repair <номер>)
func HandleRepair(conn net.Conn, cmd string, p *player.Player, npcRepo *npc_repo.PostgresNPCRepository) {
	if !p.IsTalkin || p.TalkToID != "blacksmith" {
		fmt.Fprintf(conn, "Ты не разговариваешь с кузнецом\n> ")
		return
	}

	args := strings.Fields(cmd)
	if len(args) < 2 {
		fmt.Fprintf(conn, "Что чинить? Используй `repair <номер>`\n> ")
		return
	}

	num, err := strconv.Atoi(args[1])
	if err != nil || num < 1 {
		fmt.Fprintf(conn, "Неверный номер\n> ")
		return
	}

	//получение кузнеца из БД
	target, err := npcRepo.FindByID(p.TalkToID)
	if err != nil || target == nil {
		fmt.Fprintf(conn, "Кузнец не найден\n> ")
		return
	}

	//список индексов поврежденных предметов
	var brokenItems []BrokenRef
	canRepair := map[string]bool{
		"weapon": true,
		"armor":  true,
		"helmet": true,
		"shield": true,
		"boots":  true,
		"bag":    true,
		"ring":   true,
	}
	//поиск в инвентаре
	for i, item := range p.Inventory {
		if item.Durability < 100 && canRepair[item.ItemType] {
			brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: false})
		}
	}

	//в мешке
	if p.Equipment.Bag != nil {
		for i, item := range p.Equipment.BagItems {
			if item.Durability < 100 && canRepair[item.ItemType] {
				brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: true})
			}
		}
	}

	//проверка границ
	if num > len(brokenItems) {
		fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", num)
		return
	}

	//сам предмет
	ref := brokenItems[num-1] //индекс
	var itemToRepair *item.ItemStack

	if ref.InBag {
		itemToRepair = p.Equipment.BagItems[ref.Index]
	} else {
		itemToRepair = p.Inventory[ref.Index]
	}

	price := npc.CalculateRepairPrice(itemToRepair)

	//нужен ли ремонт вообще
	if price == 0 {
		fmt.Fprintf(conn, "Этот предмет не требует ремонта\n> ")
		return
	}

	if !p.HasItemGlobal("coin", price) {
		fmt.Fprintf(conn, "Недостаточно монет! Нужно %d\n> ", price)
		return
	}

	duration := npc.CalculateRepairTime(itemToRepair)

	//создание заказа
	order := &npc.RepairOrder{
		PlayerID:  p.ID,
		Item:      itemToRepair,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(duration),
		Price:     price,
	}

	p.RemoveOneItem(itemToRepair.Name, ref.InBag, ref.Index)

	p.RemoveItemGlobal("coin", price)

	//добавление заказа кузнецу
	target.Mu.Lock()
	target.RepairOrders = append(target.RepairOrders, order)
	target.Mu.Unlock()

	npcRepo.Save(target)

	fmt.Fprintf(conn, "🔨 Кузнец взял %s в ремонт.\n> ", item.GetColoredName(itemToRepair))
	fmt.Fprintf(conn, "Готово будет через %v\n", duration)
	fmt.Fprintf(conn, "Списано: %d монет.\n> ", price)

}

// забрать готовое снаряжени (take <номер>)
func HandleTakeRepaired(cmd string, conn net.Conn, p *player.Player, npcRepo *npc_repo.PostgresNPCRepository) {
	if !p.IsTalkin || p.TalkToID != "blacksmith" {
		fmt.Fprintf(conn, "Ты не разговариваешь с кузнецом\n> ")
		return
	}

	args := strings.Fields(cmd)
	if len(args) < 2 {
		fmt.Fprintf(conn, "Что забрать? Используй `take <номер>`\n> ")
		return
	}

	num, err := strconv.Atoi(args[1])
	if err != nil || num < 1 {
		fmt.Fprintf(conn, "Неверный номер\n> ")
		return
	}

	//получение кузнеца из БД
	target, err := npcRepo.FindByID(p.TalkToID)
	if err != nil || target == nil {
		fmt.Fprintf(conn, "Кузнец не найден\n> ")
		return
	}

	var brokenItems []BrokenRef // индексы в p.Inventory
	canRepair := map[string]bool{
		"weapon": true,
		"armor":  true,
		"helmet": true,
		"shield": true,
		"boots":  true,
		"bag":    true,
		"ring":   true,
	}

	for i, item := range p.Inventory {
		if item.Durability < 100 && canRepair[item.ItemType] {
			brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: false})
		}
	}

	//в мешке
	if p.Equipment.Bag != nil {
		for i, item := range p.Equipment.BagItems {
			if item.Durability < 100 && canRepair[item.ItemType] {
				brokenItems = append(brokenItems, BrokenRef{Index: i, InBag: true})
			}
		}
	}

	//кол-во именно твоих заказов
	myOrders := []*npc.RepairOrder{}
	for _, order := range target.RepairOrders {
		if order.PlayerID == p.ID {
			myOrders = append(myOrders, order)
		}
	}

	startIdx := len(brokenItems) + 1 //первый номер заказа
	if num < startIdx || num > startIdx+len(myOrders)-1 {
		fmt.Fprintf(conn, "Неверный номер заказа\n> ")
		return
	}

	orderIdx := num - startIdx
	order := myOrders[orderIdx]

	if time.Now().Before(order.EndTime) {
		fmt.Fprintf(conn, "Заказ еще не готов\n> ")
		return
	}

	order.Item.Durability = 100

	if !p.AddItemToInventory(order.Item) {
		fmt.Fprintf(conn, "Нет места в инвентаре\n> ")
		return
	}

	//удаляем заказ у кузнеца
	target.Mu.Lock()

	for i, o := range target.RepairOrders {
		if o == order {
			target.RepairOrders = append(target.RepairOrders[:i], target.RepairOrders[i+1:]...)
			break
		}
	}

	target.Mu.Unlock()
	npcRepo.Save(target)

	fmt.Fprintf(conn, "🔨 Ты забрал %s (прочность: 100/100)\n> ", item.GetColoredName(order.Item))
}
