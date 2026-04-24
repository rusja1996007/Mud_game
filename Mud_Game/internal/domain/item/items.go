package item

// ItemData хранит базовые характеристики предмета
type ItemData struct {
	Name          string
	ItemType      string
	HungerRestore int
	ThirstRestore int
	SlotBonus     int
}

// ItemsDB — база данных всех предметов в игре
var ItemsDB = map[string]ItemData{
	//////////////////////////////////////liquid container////////////////////////////////
	"empty bottle": {
		Name:     "empty bottle",
		ItemType: "liquid container",
	},

	//////////////////////////////////////drink///////////////////////////////////////////
	"water bottle": {
		Name:          "water bottle",
		ItemType:      "drink",
		ThirstRestore: 30,
	},

	//////////////////////////////////////food////////////////////////////////////////////
	"tomato": {
		Name:          "tomato",
		ItemType:      "food",
		HungerRestore: 10,
	},

	"potato": {
		Name:          "potato",
		ItemType:      "food",
		HungerRestore: 10,
	},

	"rubus caesius": { //Ежевика
		Name:          "rubus caesius",
		ItemType:      "food",
		HungerRestore: 5,
	},

	/////////////////////////////////////container////////////////////////////////////////
	"empty bag": {
		Name:     "empty bag",
		ItemType: "container",
	},

	/////////////////////////////////////bag//////////////////////////////////////////////
	"leather bag": {
		Name:      "leather bag",
		ItemType:  "bag",
		SlotBonus: 4,
	},

	/////////////////////////////////////seed/////////////////////////////////////////////
	"tomato seeds": {
		Name:     "tomato seeds",
		ItemType: "seed",
	},

	"potato seeds": {
		Name:     "potato seeds",
		ItemType: "seed",
	},

	"burdock seeds": {
		Name:     "burdock seeds",
		ItemType: "seed",
	},

	"clover seeds": {
		Name:     "clover seeds",
		ItemType: "seed",
	},

	////////////////////////////////////weapon//////////////////////////////////////////////
	"iron sword": {
		Name:     "iron sword",
		ItemType: "weapon",
	},

	"knife": {
		Name:     "knife",
		ItemType: "weapon",
	},

	///////////////////////////////////helmet///////////////////////////////////////////////
	"leather hood": {
		Name:     "leather hood",
		ItemType: "helmet",
	},

	//////////////////////////////////armor/////////////////////////////////////////////////
	"leather armor": {
		Name:     "leather armor",
		ItemType: "armor",
	},

	//////////////////////////////////boots//////////////////////////////////////////////////
	"leather boots": {
		Name:     "leather boots",
		ItemType: "boots",
	},

	/////////////////////////////////shield//////////////////////////////////////////////////
	"wooden shield": {
		Name:     "wooden shield",
		ItemType: "shield",
	},

	/////////////////////////////////ring///////////////////////////////////////////////////
	"silver ring": {
		Name:     "silver ring",
		ItemType: "ring",
	},

	"gold ring": {
		Name:     "gold ring",
		ItemType: "ring",
	},

	"black ring": {
		Name:     "black ring",
		ItemType: "ring",
	},

	"cooper ring": { //медное
		Name:     "cooper ring",
		ItemType: "ring",
	},

	/////////////////////////////////ingredients////////////////////////////////////////////
	"burdock": { //лопух
		Name:     "burdock",
		ItemType: "ingredients",
	},

	"clover": { //клевер
		Name:     "clover",
		ItemType: "ingredients",
	},
	"inonotus obliquus": { //гриб чага
		Name:     "inonotus obliquus",
		ItemType: "ingredients",
	},
	"rubroboletus satanas": { //сатанинский гриб
		Name:     "rubroboletus satanas",
		ItemType: "ingredients",
	},
	"boletus edulis": { //белый гриб
		Name:     "boletus edulis",
		ItemType: "ingredients",
	},
	"hare meat": { //мясо зайца
		Name:     "hare meat",
		ItemType: "ingredients",
	},
	"hare ears": { //уши зайца
		Name:     "hare ears",
		ItemType: "ingredients",
	},
	"hare paws": { //лапы зайца
		Name:     "hare paws",
		ItemType: "ingredients",
	},
	"broken sword": {
		Name:     "broken sword",
		ItemType: "ingredients",
	},

	////////////////////////////////валюта//////////////////////////////////////////////
	"coin": { //монета
		Name:     "coin",
		ItemType: "currency",
	},
}

// возвращает копию предмета из базы
func GetItem(name string, count int) *ItemStack {
	data, ok := ItemsDB[name]

	if !ok {
		return nil
	}

	return &ItemStack{
		Name:          data.Name,
		Count:         count,
		ItemType:      data.ItemType,
		HungerRestore: data.HungerRestore,
		ThirstRestore: data.ThirstRestore,
		SlotBonus:     data.SlotBonus,
	}
}
