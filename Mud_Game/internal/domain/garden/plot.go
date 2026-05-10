package garden

// ✅ plot.go - что такое грядка (содержит растение или nil)
// грядки
type Plot struct {
	ID    int    // номер грядки (0, 1, 2)
	Plant *Plant // nil если пусто
}

// Сад
type Garden struct {
	PlayerID string
	Plots    []*Plot // все грядки игрока
}

func NewGarden(platerID string, plotCount int) *Garden {
	plots := make([]*Plot, plotCount)
	for i := 0; i < plotCount; i++ {
		plots[i] = &Plot{
			ID:    i,
			Plant: nil, //пустая грядка
		}
	}
	return &Garden{
		PlayerID: platerID,
		Plots:    plots,
	}
}

// посадить растение
func (g *Garden) Plant(plotID int, plantType PlantType) bool {
	if plotID < 0 || plotID >= len(g.Plots) {
		return false // нет такой грядки
	}
	if g.Plots[plotID].Plant != nil {
		return false //грядка занята
	}
	g.Plots[plotID].Plant = NewPlant(plantType)
	return true
}

// собрать урожай
// int - кол-во
// int - XP
func (g *Garden) Harvest(plotID int) (PlantType, int, bool, int) {
	if plotID < 0 || plotID >= len(g.Plots) {
		return "", 0, false, 0
	}

	plant := g.Plots[plotID].Plant
	if plant == nil || !plant.IsReady() {
		return "", 0, false, 0
	}
	yield, _ := GetPlantYield(plant.Type)

	//очищаем грядку
	g.Plots[plotID].Plant = nil

	exp := 10 //за сбор 1 растения

	return plant.Type, yield, true, exp

}
