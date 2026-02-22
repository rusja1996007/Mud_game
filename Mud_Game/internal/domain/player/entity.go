package player

type Player struct {
	ID          string
	Name        string
	CurrentRoom string   //текущая комната
	Inventory   []string // ← сюда будем складывать ID предметов
}
