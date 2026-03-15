package garden

//✅ plant.go - что такое растение

import (
	"math/rand"
	"time"
)

type PlantType string //какой то тип растения или продукта

const ( //Это как переменная, но её значение нельзя изменить после объявления.
	PlantTomato       PlantType = "Tomat"
	PlantPotato       PlantType = "Potat"
	PlantMeadowClover PlantType = "MeadowClover"
	PlantBurdock      PlantType = "Burdock"
)

// Время роста для каждого растения (в минутах)
var PlantGrowTime = map[PlantType]time.Duration{
	PlantTomato:       1 * time.Minute,
	PlantPotato:       1 * time.Minute,
	PlantMeadowClover: 2 * time.Minute,
	PlantBurdock:      2 * time.Minute,
}

// Генератор чисел на все приложение
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// делаем функцию, которая возвращает рандомный урожай
// bool = true  - всё ок, урожай есть
// bool = false - что-то пошло не так(посадили не то)
func GetPlantYield(plantType PlantType) (int, bool) {
	switch plantType { // Смотрим, какую кнопку нажал игрок
	case PlantTomato: //  Если нажал "Помидор" 🍅
		return rng.Intn(4) + 1, true // ВОЗВРАЩАЕТ число от 1 до 4
	case PlantPotato:
		return rng.Intn(4) + 1, true
	case PlantMeadowClover:
		return rng.Intn(5) + 2, true //2-6
	case PlantBurdock:
		return rng.Intn(3) + 1, true
	default:
		return 0, false
	}
}

type Plant struct {
	Type      PlantType
	PlantedAt time.Time
}

func NewPlant(plantType PlantType) *Plant {
	return &Plant{
		Type:      plantType,
		PlantedAt: time.Now(),
	}
}

func (p *Plant) IsReady() bool {
	growTime := PlantGrowTime[p.Type]     // получаем время роста
	timePassed := time.Since(p.PlantedAt) // сколько прошло времени
	return timePassed >= growTime         //// готово если прошло достаточно времени
}
