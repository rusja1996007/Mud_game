package handlers

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func HandleTake(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	//ВЗЯТЬ
	if cmd == "take" {
		fmt.Fprintf(conn, "Что взять?\n> ")
		return
	}
	args, found := strings.CutPrefix(cmd, "take ") //CutPrefix проверяет, начинается ли с "take "
	if !found {
		return // это не take, просто идём дальше к другим командам
	} //Если да → парсим команду
	args = strings.TrimSpace(args) // ✅ Обрезаем пробелы!
	if args == "" {
		fmt.Fprintf(conn, "Что взять?\n> ")
		return
	}
	parts := strings.Fields(args) //разбивает по пробелам
	//args	                  parts
	//"3 bottle"	    ["3", "bottle"]
	//"all bottle"	    ["all", "bottle"]
	//"bottle"	        ["bottle"]
	//"3 big bottle"	["3", "big", "bottle"]

	r, err := roomRepo.FindByID(p.CurrentRoom) // узнали в какой сейчас комнате
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
		return
	}

	var count int = 1 // сколько предметов брать,по умолчанию берём 1
	var itemName string

	//Парсим → count=3, itemName="Empty bottle"(пример)
	// Смотрим, что нам прислали

	if len(parts) == 1 && parts[0] == "all" { //// Это команда "take all" — взять всё из комнаты
		allItems := r.GetItems()         //получаем все предметы из комнаты
		for _, stack := range allItems { //проходим по каждому предмету
			takenStack, err := r.TakeItem(stack.Name, stack.Count) // берём сразу всю стопку
			if err != nil {
				fmt.Fprintf(conn, "Ошибка при взятии %s: %s\n> ", stack.Name, err.Error())
				continue
			}
			p.Inventory = AddToInventory(p.Inventory, takenStack)
		}

		playerRepo.Save(p)
		roomRepo.Save(r)
		fmt.Fprintf(conn, "Вы взяли все из комнаты\n> ")
		return

	} else if len(parts) == 1 {
		itemName = parts[0]
	} else {
		//Случай Б: 2 или больше частей
		if parts[0] == "all" { //Если первое слово — "all"
			count = -1                              //специальное значение "все"
			itemName = strings.Join(parts[1:], " ") //itemName = strings.Join(["big", "bottle"], " ") = "big bottle"
		} else {
			num, err := strconv.Atoi(parts[0]) // ← пробуем распарсить число
			if err == nil && num > 0 {         //значит это число
				count = num
				itemName = strings.Join(parts[1:], " ")
			} else if err == nil && num <= 0 {
				fmt.Fprintf(conn, "Количество должно быть положительным числом\n> ")
				return
			} else { // это не число и не "all" — значит, название из нескольких слов
				itemName = strings.Join(parts, " ")
			}
		}
	}
	if itemName == "" {
		fmt.Fprintf(conn, "Что взять?\n> ")
		return
	}

	//  Если это не "take all", обрабатываем обычный take
	if itemName != "" {
		items := r.GetItems()
		foundIndex := -1
		for i, stack := range items {
			if stack.Name == itemName {
				foundIndex = i
				break
			}
		}
		if foundIndex == -1 {
			fmt.Fprintf(conn, "Здесь нет такого предмета\n> ")
			return
		}
		// Смотрим, сколько штук этого предмета лежит в комнате.
		available := items[foundIndex].Count
		// Сколько будем брать
		takeCount := count
		if count == -1 {
			takeCount = available //если "all" — берём всё
		}
		if takeCount > available {
			takeCount = available // нельзя взять больше, чем есть
		}
		if takeCount == 0 {
			fmt.Fprintf(conn, "Нечего брать\n> ")
			return
		}

		//берем предметы
		takenStack, err := r.TakeItem(itemName, takeCount)
		if err != nil {
			fmt.Fprintf(conn, "Не получилось взять предметы\n> ")
			return
		}
		p.Inventory = AddToInventory(p.Inventory, takenStack)
		playerRepo.Save(p)
		roomRepo.Save(r)

		fmt.Fprintf(conn, "Ты взял %d %s\n> ", takeCount, itemName)

	}

}
