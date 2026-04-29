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

	var count int = 1 //по умолчанию
	var itemName string

	// Проверка на номер предмета
	if len(parts) == 1 {
		if num, err := strconv.Atoi(parts[0]); err == nil {
			items := r.GetItems()
			if num < 1 || num > len(items) {
				fmt.Fprintf(conn, "Нет предмета с номером %d\n> ", num)
				return
			}
			stack := items[num-1]
			itemName = stack.Name
			count = stack.Count
			goto takeItem
		} else {
			itemName = parts[0]
		}
	} else if len(parts) > 1 {
		// обычная обработка (all, число + название и т.д.)
		if parts[0] == "all" {
			count = -1
			itemName = strings.Join(parts[1:], " ")
		} else {
			num, err := strconv.Atoi(parts[0])
			if err == nil && num > 0 {
				count = num
				itemName = strings.Join(parts[1:], " ")
			} else if err == nil && num <= 0 {
				fmt.Fprintf(conn, "Количество должно быть положительным числом\n> ")
				return
			} else {
				itemName = strings.Join(parts, " ")
			}
		}
	}

	if itemName == "" {
		fmt.Fprintf(conn, "Что взять?\n> ")
		return
	}

	//Парсим → count=3, itemName="Empty bottle"(пример)
	// Смотрим, что нам прислали

	if len(parts) == 1 && parts[0] == "all" { //// Это команда "take all" — взять всё из комнаты
		// Берем предметы ПОКА они есть, а не по фиксированному списку!
		taken := 0
		stopped := false // флаг, что остановились из-за нехватки места

		for {
			items := r.GetItems()
			if len(items) == 0 {
				break
			}

			// Берем первый предмет из списка
			stack := items[0]
			takenStack, err := r.TakeItem(stack.Name, stack.Count)
			if err != nil {
				fmt.Fprintf(conn, "Ошибка при взятии предмета: %s", err.Error())
				continue
			}

			if !p.AddItemToInventory(takenStack) {
				// Если нет места - кладем предмет обратно в комнату
				r.AddItem(takenStack)
				fmt.Fprintf(conn, "Нет места в инвентаре! Остановлено на %s\n", takenStack.Name)
				stopped = true
				break
			}
			taken++

		}
		if taken > 0 {
			playerRepo.Save(p)
			roomRepo.Save(r)
			if !stopped {
				fmt.Fprintf(conn, "Вы взяли все из комнаты\n> ")
			} else {
				fmt.Fprintf(conn, "Вы взяли %d предметов, Освободите место для остальных.\n> ", taken)
			}
		} else {
			fmt.Fprintf(conn, "В комнате нет предметов\n> ")
		}
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

takeItem:
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
		if !p.AddItemToInventory(takenStack) {
			r.AddItem(takenStack)
			fmt.Fprintf(conn, "Нет места в инвентаре!\n> ")
			return
		}
		playerRepo.Save(p)
		roomRepo.Save(r)

		fmt.Fprintf(conn, "Ты взял %d %s\n> ", takeCount, itemName)

	}

}
