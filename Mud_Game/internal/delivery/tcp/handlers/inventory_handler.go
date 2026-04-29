package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
)

func HandleInventory(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//ИНВЕНТАРЬ
	var result strings.Builder
	//создаем Equipment если его нет
	if p.Equipment == nil {
		p.Equipment = &player.Equipment{}
		playerRepo.Save(p)
	}

	//рюкзак:
	fmt.Fprintf(&result, "Инвентарь: %d/%d слотов\n", p.GetUsedSlots(), p.GetMaxSlots())

	if len(p.Inventory) == 0 {
		fmt.Fprintf(&result, "Инвентарь пуст\n> ")
	} else {
		fmt.Fprintf(&result, "\nПредметы в Инвентаре:\n")
		for i, stack := range p.Inventory {
			fmt.Fprintf(&result, " %d. %s", i+1, stack.Name)
			if stack.Durability > 0 {
				fmt.Fprintf(&result, " (прочность: %d)", stack.Durability)
			}
			if stack.Count > 1 {
				fmt.Fprintf(&result, " x%d", stack.Count)
			}
			fmt.Fprintf(&result, "\n")
		}
	}
	//Экипировка
	fmt.Fprintf(&result, "\n Экипировка:\n")
	//проверяем каждый слот экипировки:
	//оружие
	if p.Equipment.Weapon != nil {
		fmt.Fprintf(&result, " Оружие : %s\n", p.Equipment.Weapon.Name)
	} else {
		fmt.Fprintf(&result, " Оружие : не надето\n")
	}
	//броня
	if p.Equipment.Armor != nil {
		fmt.Fprintf(&result, " Броня : %s\n", p.Equipment.Armor.Name)
	} else {
		fmt.Fprintf(&result, " Броня : не надето\n")
	}

	//шлем
	if p.Equipment.Helmet != nil {
		fmt.Fprintf(&result, " Шлем : %s\n", p.Equipment.Helmet.Name)
	} else {
		fmt.Fprintf(&result, " Шлем : не надето\n")
	}

	//мешок
	if p.Equipment.Bag != nil {
		fmt.Fprintf(&result, " Мешок : %s\n", p.Equipment.Bag.Name)
	} else {
		fmt.Fprintf(&result, " Мешок : не надето\n")
	}

	//Щит
	if p.Equipment.Shield != nil {
		fmt.Fprintf(&result, " Щит : %s\n", p.Equipment.Shield.Name)
	} else {
		fmt.Fprintf(&result, " Щит : не надето\n")
	}

	//обувь
	if p.Equipment.Boots != nil {
		fmt.Fprintf(&result, " Обувь : %s\n", p.Equipment.Boots.Name)
	} else {
		fmt.Fprintf(&result, " Обувь : не надето\n")
	}

	//кольцо1
	if p.Equipment.Ring1 != nil {
		fmt.Fprintf(&result, " Кольцо_1 : %s\n", p.Equipment.Ring1.Name)
	} else {
		fmt.Fprintf(&result, " Кольцо_1 : не надето\n")
	}

	//кольцо2
	if p.Equipment.Ring2 != nil {
		fmt.Fprintf(&result, " Кольцо_2 : %s\n", p.Equipment.Ring2.Name)
	} else {
		fmt.Fprintf(&result, " Кольцо_2 : не надето\n")
	}
	fmt.Fprintf(&result, "> ")
	conn.Write([]byte(result.String()))

}
