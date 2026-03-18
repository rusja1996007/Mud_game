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

// собрать содержимое грядок(подсказки в plant_handler)
func HandleHarvest(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	args := strings.Fields(cmd)

	if p.CurrentRoom != p.Zone.GardenID {
		fmt.Fprintf(conn, "Ты можешь собирать урожай только в своем огороде!\n> ")
		return
	}

	if len(args) < 2 {
		fmt.Fprintf(conn, "Использование: harvest <номер_грядки>\n> ")
		return
	}

	plotID, err := strconv.Atoi(args[1])
	if err != nil || plotID > 3 || plotID < 1 {
		fmt.Fprintf(conn, "Номер грядки должен быть 1,2 или 3!\n> ")
		return
	}

	plotID--

	plantType, vsego, ok := p.Zone.Garden.Harvest(plotID)
	if !ok {
		fmt.Fprintf(conn, "На этой грядке ничего нет для сбора или растение еще не выросло!\n> ")
		return
	}

	var itemName string
	switch plantType {
	case garden.PlantTomato:
		itemName = "tomato"
	case garden.PlantPotato:
		itemName = "potato"
	case garden.PlantBurdock:
		itemName = "burdock"
	case garden.PlantMeadowClover:
		itemName = "meadow_clover"
	}

	player.AddItem(&p.Inventory, itemName, vsego)

	playerRepo.Save(p)

	fmt.Fprintf(conn, "Ты собрал %d %s с грядки %d\n> ", vsego, itemName, plotID+1) //// +1 чтобы показать игроку 1,2,3 а не 0,1,2

}
