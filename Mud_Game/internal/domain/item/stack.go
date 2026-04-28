package item

// структуре представляет стопку одинаковых предметов
type ItemStack struct {
	Name      string `json:"name"`
	Count     int    `json:"count"`      //сколько штук
	SlotBonus int    `json:"slot_bonus"` //сколько слотов дает

	ItemType      string `json:"item_type"`      //тип предмета, в какую ячейку можно отнести(
	HungerRestore int    `json:"hunger_restore"` //сколько восстанавливает еды
	ThirstRestore int    `json:"thirst_restore"` //сколько восстанавливает жажды

	MinDamage  int `json:"min_damage"`
	MaxDamage  int `json:"max_damage"`
	Durability int `json:"durability"` //прочность
}

/*
weapon    - оружие
armor     - броня
helmet    - шлем
shield    - щит
boots     - обувь
ring      - кольцо
bag       - мешок (расширяет слоты)
food      - еда
drink     - напиток
seed      - семена
container - контейнер (мешок)
liquid container - бутылка
material  - материал для крафта
currency  - валюта
*/

// уменьшение прочности
func (i *ItemStack) Decrease(amount int) bool {
	i.Durability -= amount
	if i.Durability <= 0 {
		return true //сломался
	}
	return false
}
