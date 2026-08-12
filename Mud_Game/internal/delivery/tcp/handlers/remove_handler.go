package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
)

// Снятие экипировки
func HandleRemove(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	args, found := strings.CutPrefix(cmd, "remove ")
	if !found {
		return
	}
	args = strings.TrimSpace(args)
	if args == "" {
		fmt.Fprintf(conn, "Что снять? Использование : remove <предмет>\n> ")
		return
	}

	//нижний регистр
	slot := strings.ToLower(args)

	//в зависимотсти от слота:
	switch slot {
	case "weapon":
		if p.Equipment.Weapon != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Weapon
			p.AddItemToInventory(p.Equipment.Weapon)
			p.Equipment.Weapon = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надето оружие\n> ")
		}
	case "armor":
		if p.Equipment.Armor != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Armor
			p.AddItemToInventory(p.Equipment.Armor)
			p.Equipment.Armor = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надета броня\n> ")
		}
	case "helmet":
		if p.Equipment.Helmet != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Helmet
			p.AddItemToInventory(p.Equipment.Helmet)
			p.Equipment.Helmet = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надет шлем\n> ")
		}
	case "bag":
		if p.Equipment.Bag != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Bag
			p.AddItemToInventory(p.Equipment.Bag)
			p.Equipment.Bag = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надета сумка\n> ")
		}
	case "shield":
		if p.Equipment.Shield != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Shield
			p.AddItemToInventory(p.Equipment.Shield)
			p.Equipment.Shield = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надет щит\n> ")
		}
	case "boots":
		if p.Equipment.Boots != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Boots
			p.AddItemToInventory(p.Equipment.Boots)
			p.Equipment.Boots = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надета обувь\n> ")
		}
	case "ring1":
		if p.Equipment.Ring1 != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Ring1
			p.AddItemToInventory(p.Equipment.Ring1)
			p.Equipment.Ring1 = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надето кольцо на первой руке\n> ")
		}
	case "ring2":
		if p.Equipment.Ring2 != nil {
			if !p.CanAddItem() {
				fmt.Fprintf(conn, "Нет места в инвентаре! Освободи место сначала.\n> ")
				return
			}
			itemName := p.Equipment.Ring2
			p.AddItemToInventory(p.Equipment.Ring2)
			p.Equipment.Ring2 = nil
			fmt.Fprintf(conn, "Ты снял %s\n> ", item.GetColoredName(itemName))
		} else {
			fmt.Fprintf(conn, "У тебя не надето кольцо на второй руке\n> ")
		}
	default:
		fmt.Fprintf(conn, "Неизвестный слот\n> ")
		return
	}
	playerRepo.Save(p)
}
