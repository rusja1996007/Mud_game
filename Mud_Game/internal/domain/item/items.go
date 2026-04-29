package item

// ItemData хранит базовые характеристики предмета
type ItemData struct {
	Name          string
	ItemType      string
	HungerRestore int
	ThirstRestore int
	SlotBonus     int

	MinDamage int
	MaxDamage int

	Durability int //прочность
	Defence    int //защита

	Description string
	ID          string
}

// ItemsDB — база данных всех предметов в игре
var ItemsDB = map[string]ItemData{
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

	//////////////////////////////////////food////////////////////////////////////////////
	"tomato": {
		Name:          "tomato",
		ItemType:      "food",
		HungerRestore: 10,
		Description:   "--",
	},

	"potato": {
		Name:          "potato",
		ItemType:      "food",
		HungerRestore: 10,
		Description:   "--",
	},

	"rubus caesius": { //Ежевика
		Name:          "rubus caesius",
		ItemType:      "food",
		HungerRestore: 5,
		Description:   "--",
	},

	/////////////////////////////////////container////////////////////////////////////////
	"empty bag": {
		Name:        "empty bag",
		ItemType:    "container",
		Description: "--",
	},

	/////////////////////////////////////bag//////////////////////////////////////////////
	"leather bag": {
		Name:       "leather bag",
		ItemType:   "bag",
		SlotBonus:  4,
		Durability: 100,
		Defence:    1,
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
		Name:       "iron sword",
		ItemType:   "weapon",
		MinDamage:  20,
		MaxDamage:  30,
		Durability: 100,
	},

	"knife": {
		Name:       "knife",
		ItemType:   "weapon",
		MinDamage:  5,
		MaxDamage:  10,
		Durability: 100,
	},

	///////////////////////////////////helmet///////////////////////////////////////////////
	"leather hood": {
		Name:       "leather hood",
		ItemType:   "helmet",
		Durability: 100,
		Defence:    6,
	},

	//////////////////////////////////armor/////////////////////////////////////////////////
	"leather armor": {
		Name:       "leather armor",
		ItemType:   "armor",
		Durability: 100,
		Defence:    10,
	},

	//////////////////////////////////boots//////////////////////////////////////////////////
	"leather boots": {
		Name:       "leather boots",
		ItemType:   "boots",
		Durability: 100,
		Defence:    5,
	},

	/////////////////////////////////shield//////////////////////////////////////////////////
	"wooden shield": {
		Name:       "wooden shield",
		ItemType:   "shield",
		Durability: 100,
		Defence:    20,
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
		Description: "--",
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
		Description: "--",
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
	}
}
