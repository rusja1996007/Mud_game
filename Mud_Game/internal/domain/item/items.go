package item

// ItemData хранит базовые характеристики предмета
type ItemData struct {
	Name          string
	ItemType      string
	HungerRestore int
	ThirstRestore int
	SlotBonus     int
	HungerRate    int
	ThirstRate    int
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

	/////////////////////////////////ingredients////////////////////////////////////////////
	"burdock": {
		Name:     "burdock",
		ItemType: "ingredients",
	},

	"clover": {
		Name:     "clover",
		ItemType: "ingredients",
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
		HungerRate:    data.HungerRate,
		ThirstRate:    data.ThirstRate,
	}
}
