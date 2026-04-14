package handlers

import (
	"Mud_game/Mud_Game/internal/domain/garden"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
	"time"
)

// просмотр огорода
func HandleGarden(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	if p.Zone == nil {
		fmt.Fprintf(conn, "У тебя нет личной зоны! ОБратись к админу!\n> ")
		return
	}

	if p.Zone.Garden == nil {
		fmt.Fprintf(conn, "у тебя нет огорода\n> ")
		return
	}

	if p.CurrentRoom != p.Zone.GardenID {
		fmt.Fprintf(conn, "Тебе надо быть в своем огороде чтобы осмотреть его!\n> ")
		return
	}

	gardenObj := p.Zone.Garden
	if gardenObj == nil {
		fmt.Fprintf(conn, "У тебя нет огорода!\n")
		return
	}

	//Создание "буфера" для ответа:
	var result strings.Builder
	result.WriteString("🌱 Твой огород:\n")

	//Проходим по всем грядкам
	for i, plot := range gardenObj.Plots {
		gardenNumber := i + 1
		//Проверь, пустая ли грядка
		if plot.Plant == nil {
			fmt.Fprintf(&result, "Грядка %d: пусто!\n", gardenNumber)
			continue
		}
		//Если есть растение:
		plant := plot.Plant

		//Достаем информацию о растении.
		var plantName string
		switch plant.Type {
		case garden.PlantTomato:
			plantName = "tomato"
		case garden.PlantPotato:
			plantName = "potato"
		case garden.PlantBurdock:
			plantName = "burdock"
		case garden.PlantClover:
			plantName = "clover"
		}

		// если готов-
		if plant.IsReady() {
			fmt.Fprintf(&result, "Грядка %d: %s 🟢 ГОТОВО!\n", gardenNumber, plantName)
		} else {
			//если нет- сколько осталось

			// Сколько всего нужно расти
			growTime := garden.PlantGrowTime[plant.Type]
			//Спрашиваем у компьютера: "Сколько времени прошло с момента посадки?"
			timePassed := time.Since(plant.PlantedAt)
			//Из общего времени роста вычитаем уже прошедшее время.
			timeLeft := growTime - timePassed
			//"Дай это время в минутах, просто числом".
			minutes := int(timeLeft.Minutes())
			//Если растение уже созрело, timeLeft может быть отрицательным. Мы говорим: "Если меньше 0, покажи 0"(защита)
			minutes = max(minutes, 0)

			fmt.Fprintf(&result, "Грядка %d: %s ⏳ растет, осталось %d минут!\n", gardenNumber, plantName, minutes)
		}
	}
	fmt.Fprintf(&result, "\nКоманды: Посадить-plant <номер> <растение>,Собрать-harvest <номер>\n> ")

	conn.Write([]byte(result.String()))

}
