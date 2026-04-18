package item

// структуре представляет стопку одинаковых предметов
type ItemStack struct {
	Name      string `json:"name"`
	Count     int    `json:"count"`      //сколько штук
	SlotBonus int    `json:"slot_bonus"` //сколько слотов дает

	ItemType      string `json:"item_type"`      //тип предмета, в какую ячейку можно отнести(
	HungerRestore int    `json:"hunger_restore"` //сколько восстанавливает еды
	ThirstRestore int    `json:"thirst_restore"` //сколько восстанавливает жажды

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
container - контейнер (бутылка, мешок)
material  - материал для крафта
*/
