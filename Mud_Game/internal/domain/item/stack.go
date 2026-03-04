package item

// структуре представляет стопку одинаковых предметов
type ItemStack struct {
	Name  string `json:"name"`
	Count int    `json:"count"` //сколько штук
}
