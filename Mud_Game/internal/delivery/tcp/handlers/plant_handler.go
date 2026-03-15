package handlers

import (
	"Mud_game/Mud_Game/internal/domain/garden"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func HandlePlant(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {

	args := strings.Fields(cmd) // - разбивает строку на слова, убирает лишние пробелы вначале и конце

	// 1. проверяем что игрок в огороде
	if p.CurrentRoom != p.Zone.GardenID {
		fmt.Fprintf(conn, "Ты можешь сажать только в своем огороде!\n> ")
		return
	}

	// 2. проверяем аргументы
	if len(args) < 3 {
		fmt.Fprintf(conn, "Использование: plant <номер грядки> <растение>\n> ")
		return
	}

	// 3. получаем номер грядки
	plotID, err := strconv.Atoi(args[1])
	if err != nil || plotID < 1 || plotID > 3 {
		fmt.Fprintf(conn, "Номер грядки должен быть 1, 2 или 3!\n> ")
		return
	}
	plotID-- //// конвертируем в индекс (0, 1, 2)
	// 4. получаем тип растения
	plantName := strings.ToLower(args[2]) //----------------------------------------
	var plantType garden.PlantType
	var seedName string //

	switch plantName {
	case "tomato":
		plantType = garden.PlantTomato
		seedName = "tomato seeds"
	case "potato":
		plantType = garden.PlantPotato
		seedName = "potato seeds"
	case "meadow_clover":
		plantType = garden.PlantMeadowClover
		seedName = "meadow_clover seeds"
	case "burdock":
		plantType = garden.PlantBurdock
		seedName = "burdock seeds"
	default:
		fmt.Fprintf(conn, "Неизвестное растение\n> ")
		return
	}

	// 5. проверяемм наличие семян
	if !player.HasItem(p.Inventory, seedName, 1) {
		fmt.Fprintf(conn, "У тебя нет %s\n> ", seedName)
		return
	}

	// 6. Пробуем посадить
	if !p.Zone.Garden.Plant(plotID, plantType) {
		fmt.Fprintf(conn, "Эта грядка занята или не существует\n> ")
		return
	}

	// 7. Убираем семена
	player.RemoveItem(&p.Inventory, seedName, 1)

	// 8. Сохраняем изменения
	playerRepo.Save(p)

	// 9. Отвечаем игроку
	fmt.Fprintf(conn, "Вы посадили %s на грядке %d\n> ", plantName, plotID+1)
}
