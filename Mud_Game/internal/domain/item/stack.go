package item

import (
	"fmt"
	"time"
)

// структуре представляет стопку одинаковых предметов
type ItemStack struct {
	Name      string `json:"name"`
	Count     int    `json:"count"`      //сколько штук
	SlotBonus int    `json:"slot_bonus"` //сколько слотов дает

	ItemType      string `json:"item_type"`      //тип предмета, в какую ячейку можно отнести(
	HungerRestore int    `json:"hunger_restore"` //сколько восстанавливает еды
	ThirstRestore int    `json:"thirst_restore"` //сколько восстанавливает жажды

	MinDamage     int `json:"min_damage"`
	MaxDamage     int `json:"max_damage"`
	HealMax       int `json:"max_heal"`
	HealMin       int `json:"min_heal"`
	Durability    int `json:"durability"` //прочность
	Defence       int `json:"defence"`    //защита
	MagicDefence  int `json:"magic_defence"`
	FireDefence   int `json:"fire_defence"`
	PoisonDefence int `json:"poison_defence"`
	MagicDamage   int `json:"magic_damage"` //урон
	FireDamage    int `json:"fire_damage"`
	PoisonDamage  int `json:"poison_damage"`

	Description string `json:"description"` //описание
	ID          string `json:"id"`          //Универсльный id

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

func GenerateItemID() string {
	return fmt.Sprintf("item_%d", time.Now().UnixNano())
}

// CanStackWith проверяет, можно ли объединить этот предмет с другим
func (i *ItemStack) CanStackWith(other *ItemStack) bool {
	if i.Name != other.Name {
		return false
	}

	switch i.ItemType {
	case "weapon", "armor", "helmet", "shield", "boots", "ring", "bag":
		return false
	}

	//остальное стакается
	return true
}
