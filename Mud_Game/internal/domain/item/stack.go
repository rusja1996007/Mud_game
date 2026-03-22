package item

// структуре представляет стопку одинаковых предметов
type ItemStack struct {
	Name       string `json:"name"`
	Count      int    `json:"count"` //сколько штук
	SlotBonus  int    //сколько слотов дает
	HungerRate int    //сколько тратит голода(для тяжелых предметов)
	ThirstRate int    //скольбко тратит жажды(для тяжелых предметов)
	ItemType   string `json:"item_type"` //тип предмета, в какую ячейку можно отнести(
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
