package item

import (
	"math/rand"
)

// ItemData хранит базовые характеристики предмета
type ItemData struct {
	Name          string
	ItemType      string
	HungerRestore int //пополнение еды
	ThirstRestore int //пополнение воды
	SlotBonus     int

	MinDamage  int
	MaxDamage  int
	DamageType string

	HealMin int
	HealMax int

	Durability    int //прочность
	Defence       int //защита(физ)
	MagicDefence  int
	FireDefence   int
	PoisonDefence int

	MagicDamage  int //урон доп на оружии
	FireDamage   int
	PoisonDamage int

	Description string
	ID          string
}

// ItemsDB — база данных всех предметов в игре
var ItemsDB = map[string]ItemData{

	//////////////////////////////////////scroll/////////////////////////////////////////
	"scroll fireball": {
		Name:        "scroll fireball",
		ItemType:    "scroll",
		Description: "Одноразовый свиток, выпускающий огненный шар (42-85 урона).",
		MinDamage:   40,
		MaxDamage:   80,
		DamageType:  "fire",
	},
	"scroll heal": {
		Name:        "scroll heal",
		ItemType:    "scroll",
		Description: "Одноразовый свиток, восстанавливающий здоровье (40-50 HP)",
		HealMin:     40,
		HealMax:     50,
	},
	//////////////////////////////////////liquid container////////////////////////////////
	"empty bottle": {
		Name:        "empty bottle",
		ItemType:    "liquid container",
		Description: "--",
	},

	//////////////////////////////////////drink///////////////////////////////////////////
	"water bottle": {
		Name:          "water bottle",
		ItemType:      "drink",
		ThirstRestore: 30,
		Description:   "--",
	},

	"antidote": {
		Name:        "antidote",
		ItemType:    "drink",
		Description: "Отвар, снимающий большинство ядов",
	},

	//////////////////////////////////////food////////////////////////////////////////////
	"tomato": {
		Name:          "tomato",
		ItemType:      "food",
		HungerRestore: 5,
		Description:   "--",
	},

	"potato": {
		Name:          "potato",
		ItemType:      "food",
		HungerRestore: 5,
		Description:   "--",
	},

	"rubus caesius": { //Ежевика
		Name:          "rubus caesius",
		ItemType:      "food",
		HungerRestore: 5,
		Description:   "--",
	},

	"vegetable set": {
		Name:          "vegetable set",
		ItemType:      "food",
		HungerRestore: 30,
		Description:   "+5 HP\n+1 HP /30 сек\nПростенький овощной набор",
	},

	/////////////////////////////////////container////////////////////////////////////////
	"empty bag": {
		Name:        "empty bag",
		ItemType:    "container",
		Description: "--",
	},

	/////////////////////////////////////bag//////////////////////////////////////////////
	"leather bag": {
		Name:        "leather bag",
		ItemType:    "bag",
		SlotBonus:   4,
		Durability:  100,
		Defence:     1,
		Description: "Кожаный мешок, позволяет носить больше предметов",
	},

	/////////////////////////////////////seed/////////////////////////////////////////////
	"tomato seeds": {
		Name:        "tomato seeds",
		ItemType:    "seed",
		Description: "--",
	},

	"potato seeds": {
		Name:        "potato seeds",
		ItemType:    "seed",
		Description: "--",
	},

	"burdock seeds": {
		Name:        "burdock seeds",
		ItemType:    "seed",
		Description: "--",
	},

	"clover seeds": {
		Name:        "clover seeds",
		ItemType:    "seed",
		Description: "--",
	},

	////////////////////////////////////weapon//////////////////////////////////////////////
	"iron sword": {
		Name:        "iron sword",
		ItemType:    "weapon",
		MinDamage:   20,
		MaxDamage:   30,
		Durability:  100,
		Description: "Обычный железный меч",
	},

	"knife": {
		Name:        "knife",
		ItemType:    "weapon",
		MinDamage:   5,
		MaxDamage:   10,
		Durability:  100,
		Description: "Обычный охотничий нож",
	},

	///////////////////////////////////helmet///////////////////////////////////////////////
	"leather hood": {
		Name:        "leather hood",
		ItemType:    "helmet",
		Durability:  100,
		Defence:     2,
		Description: "Голову от дубины спасёт. От меча — вряд ли. Но голова целее, чем без него",
	},

	//////////////////////////////////armor/////////////////////////////////////////////////
	"leather armor": {
		Name:        "leather armor",
		ItemType:    "armor",
		Durability:  100,
		Defence:     10,
		Description: "Грубая дублёная кожа — не парадная броня, но в бою не подведёт и бежать не помешает",
	},

	//////////////////////////////////boots//////////////////////////////////////////////////
	"leather boots": {
		Name:        "leather boots",
		ItemType:    "boots",
		Durability:  100,
		Defence:     3,
		Description: "Мягче железа, крепче тряпья. Шаг слышен, но не лязгает.",
	},

	/////////////////////////////////shield//////////////////////////////////////////////////
	"wooden shield": {
		Name:        "wooden shield",
		ItemType:    "shield",
		Durability:  100,
		Defence:     20,
		Description: "Лёгкий, дешёвый, шумит мало. Горит — быстро.",
	},

	/////////////////////////////////ring///////////////////////////////////////////////////
	"silver ring": {
		Name:        "silver ring",
		ItemType:    "ring",
		Description: "обычное серебряное кольцо",
	},

	"gold ring": {
		Name:        "gold ring",
		ItemType:    "ring",
		Description: "обычное золотое кольцо",
	},

	"black ring": {
		Name:        "black ring",
		ItemType:    "ring",
		Description: "черное титановое кольцо",
	},

	"cooper ring": { //медное
		Name:        "cooper ring",
		ItemType:    "ring",
		Description: "обычное медное кольцо",
	},

	/////////////////////////////////ingredients////////////////////////////////////////////
	"burdock": { //лопух
		Name:        "burdock",
		ItemType:    "ingredients",
		Description: "--",
	},

	"clover": { //клевер
		Name:        "clover",
		ItemType:    "ingredients",
		Description: "Клевер, ростет везде.",
	},
	"inonotus obliquus": { //гриб чага
		Name:        "inonotus obliquus",
		ItemType:    "ingredients",
		Description: "--",
	},
	"rubroboletus satanas": { //сатанинский гриб
		Name:        "rubroboletus satanas",
		ItemType:    "ingredients",
		Description: "--",
	},
	"boletus edulis": { //белый гриб
		Name:        "boletus edulis",
		ItemType:    "ingredients",
		Description: "Обычный белый гриб. Сырой несьедобный.",
	},
	"hare meat": { //мясо зайца
		Name:        "hare meat",
		ItemType:    "ingredients",
		Description: "--",
	},
	"hare ears": { //уши зайца
		Name:        "hare ears",
		ItemType:    "ingredients",
		Description: "--",
	},
	"hare paws": { //лапы зайца
		Name:        "hare paws",
		ItemType:    "ingredients",
		Description: "--",
	},
	"broken sword": {
		Name:        "broken sword",
		ItemType:    "ingredients",
		Description: "--",
	},

	////////////////////////////////валюта//////////////////////////////////////////////
	"coin": { //монета
		Name:        "coin",
		ItemType:    "currency",
		Description: "Повсемирная волюта",
	},
}

// возвращает копию предмета из базы
func GetItem(name string, count int) *ItemStack {
	data, ok := ItemsDB[name]

	if !ok {
		return nil
	}

	return &ItemStack{
		ID:            GenerateItemID(),
		Name:          data.Name,
		Count:         count,
		ItemType:      data.ItemType,
		HungerRestore: data.HungerRestore,
		ThirstRestore: data.ThirstRestore,
		SlotBonus:     data.SlotBonus,
		MinDamage:     data.MinDamage,
		MaxDamage:     data.MaxDamage,
		Durability:    data.Durability,
		Defence:       data.Defence,
		Description:   data.Description,
		MagicDefence:  data.MagicDefence,
		PoisonDefence: data.PoisonDefence,
		FireDefence:   data.FireDefence,
		HealMin:       data.HealMin,
		HealMax:       data.HealMax,
		MagicDamage:   data.MagicDamage,
		FireDamage:    data.FireDamage,
		PoisonDamage:  data.PoisonDamage,
	}
}

// создание + рандомный бонус для обычного меча
func CreateRandomSword() *ItemStack {
	type Bonus struct {
		Name  string
		Value int
	}

	//список бонусов (возможных)
	bonuses := []Bonus{
		{"FireDamage", 2 + rand.Intn(5)},
		{"MagicDamage", 2 + rand.Intn(5)},
		{"PoisonDamage", 2 + rand.Intn(5)},
	}

	//выбор случайного
	bonus := bonuses[rand.Intn(len(bonuses))]

	sword := &ItemStack{
		ID:          GenerateItemID(),
		Name:        "MIR",
		Count:       1,
		ItemType:    "weapon",
		MinDamage:   20 + (rand.Intn(8) + 1),
		MaxDamage:   30 + (rand.Intn(8) + 1),
		Durability:  100,
		Description: "Железный меч от которого исходит странная сила",
	}

	//применить бонус
	switch bonus.Name {
	case "FireDamage":
		sword.FireDamage = bonus.Value
	case "MagicDamage":
		sword.MagicDamage = bonus.Value
	case "PoisonDamage":
		sword.PoisonDamage = bonus.Value
	}

	return sword
}
