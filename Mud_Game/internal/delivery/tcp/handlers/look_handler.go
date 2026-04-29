package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
)

func HandleLook(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	args := strings.Fields(cmd)

	if len(args) > 1 {
		itemName := strings.Join(args[1:], " ")
		inspectItem(conn, p, roomRepo, itemName)
		return
	}

	currentRoomID := p.CurrentRoom                // 1. Получить текущую комнату игрока(id комнаты)
	room, err := roomRepo.FindByID(currentRoomID) //Обращаемся к репозиторию комнат (s.roomRepo) и просим найти комнату по этому ID.
	if err != nil {
		fmt.Fprintf(conn, "Комната  не найдена\n> ")
		return
	}
	responce := room.Look(p.ID)           //Получаем описание комнаты
	fmt.Fprintf(conn, "%s\n> ", responce) //Добавляем приглашение \n> для следующей команды
}

// ИЩет предмет и выводит его характеристики
func inspectItem(conn net.Conn, p *player.Player, roomRepo room.Repository, itemName string) {

	//поиск в инвентаре
	for _, stack := range p.Inventory {
		if stack.Name == itemName {
			printItemInfo(conn, stack)
			return
		}
	}

	//поиск в экипировке
	if p.Equipment.Weapon != nil && p.Equipment.Weapon.Name == itemName {
		printItemInfo(conn, p.Equipment.Weapon)
		return
	}

	if p.Equipment.Armor != nil && p.Equipment.Armor.Name == itemName {
		printItemInfo(conn, p.Equipment.Armor)
		return
	}
	if p.Equipment.Helmet != nil && p.Equipment.Helmet.Name == itemName {
		printItemInfo(conn, p.Equipment.Helmet)
		return
	}
	if p.Equipment.Shield != nil && p.Equipment.Shield.Name == itemName {
		printItemInfo(conn, p.Equipment.Shield)
		return
	}
	if p.Equipment.Boots != nil && p.Equipment.Boots.Name == itemName {
		printItemInfo(conn, p.Equipment.Boots)
		return
	}
	if p.Equipment.Bag != nil && p.Equipment.Bag.Name == itemName {
		printItemInfo(conn, p.Equipment.Bag)
		return
	}
	if p.Equipment.Ring1 != nil && p.Equipment.Ring1.Name == itemName {
		printItemInfo(conn, p.Equipment.Ring1)
		return
	}
	if p.Equipment.Ring2 != nil && p.Equipment.Ring2.Name == itemName {
		printItemInfo(conn, p.Equipment.Ring2)
		return
	}

	//поиск в комнате
	room, _ := roomRepo.FindByID(p.CurrentRoom)
	for _, stack := range room.GetItems() {
		if stack.Name == itemName {
			printItemInfo(conn, stack)
			return
		}
	}
	fmt.Fprintf(conn, "Ты не видишь здесь '%s'\n> ", itemName)

}

func printItemInfo(conn net.Conn, item *item.ItemStack) {
	fmt.Fprintf(conn, "📖 %s\n", item.Name)
	fmt.Fprintf(conn, "Тип %s\n", item.ItemType)

	//для оружия
	if item.ItemType == "weapon" {
		fmt.Fprintf(conn, "Урон %d-%d\n", item.MinDamage, item.MaxDamage)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	//броня
	if item.ItemType == "armor" {
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	if item.ItemType == "helmet" {
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	if item.ItemType == "shield" {
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	if item.ItemType == "boots" {
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	if item.ItemType == "bag" {
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
	}

	if item.ItemType == "ring" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "liquid container" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "drink" {
		fmt.Fprintf(conn, "Восстанавливает жажду +%d\n", item.ThirstRestore)
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "food" {
		fmt.Fprintf(conn, "Восстанавливает голод +%d\n", item.HungerRestore)
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "container" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "bag" {
		fmt.Fprintf(conn, "Слотов в инвентаре +%d\n", item.SlotBonus)
		fmt.Fprintf(conn, "Защита %d\n", item.Defence)
		fmt.Fprintf(conn, "Прочность %d/%d\n", item.Durability, 100)
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "seed" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "ingredients" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	if item.ItemType == "coin" {
		fmt.Fprintf(conn, "%s\n", item.Description)
	}

	fmt.Fprintf(conn, "> ")

}
