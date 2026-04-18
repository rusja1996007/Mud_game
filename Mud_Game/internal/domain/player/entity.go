package player

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"fmt"
	"net"
	"time"
)

// // Equipment - экипировка (не занимает слоты рюкзака)
type Equipment struct {
	Weapon *item.ItemStack
	Armor  *item.ItemStack
	Helmet *item.ItemStack
	Bag    *item.ItemStack //мешок(сумка)
	Shield *item.ItemStack //щит
	Boots  *item.ItemStack //тапки
	Ring1  *item.ItemStack //кольцо1
	Ring2  *item.ItemStack //кольцо2

}

type Player struct {
	ID          string
	Name        string
	CurrentRoom string            //текущая комната
	Inventory   []*item.ItemStack // ← стопки предметов (название + кол-во)
	Equipment   *Equipment
	Zone        *PLayerZone
	Stats       *Stats // характеристики

}

type Stats struct {
	//БАЗОВЫЕ
	MaxSlots int //слоты все в рюкзаке
	Hunger   int //голод
	Thirst   int //жажда
	Health   int // здоровье

	//характеристики
	Strength  int //сила
	Dexterity int //ловкость
	Intelect  int //интелект
	Tracking  int //следопытство

	//прогресс(позже)
	Level      int
	Experience int
}

// запускает таймер голода (каждые GetHUNGERInterval секунд )
func (p *Player) StartHungerTicker(conn net.Conn, repo Repository) {
	for {
		interval := p.GetHungerInterval()
		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		<-ticker.C //ждем один тик
		ticker.Stop()

		if p.Stats.Hunger > 0 {
			p.Stats.Hunger--
		}

		if p.Stats.Hunger == 0 {
			p.Stats.Health -= 2
			fmt.Fprintf(conn, "Ты умираешь от голода!\n> ")
		}
		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб от голода, персонаж удаляется!\n> ")
			repo.Delete(p.ID)
			conn.Close()
			return
		}
		repo.Save(p) // сохраняем только живых
	}
}

// таймер жажды (каждые getthirstticker секунд -1)
func (p *Player) StartThirstTicker(conn net.Conn, repo Repository) {
	for {
		interval := p.GetThirstInterval()
		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		<-ticker.C
		ticker.Stop()

		if p.Stats.Thirst > 0 {
			p.Stats.Thirst--
		}

		if p.Stats.Thirst == 0 {
			p.Stats.Health -= 2
			fmt.Fprintf(conn, "Ты умираешь от жажды!\n> ")
		}
		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб от обезводивания, персонаж удаляется!\n> ")
			repo.Delete(p.ID)
			conn.Close()
			return
		}
		repo.Save(p) // сохраняем только живых
	}
}

// возвращает интервал в секундах между тиками голода
func (p *Player) GetHungerInterval() int {
	interval := 60 //базово 60сек

	if p.Equipment.Bag != nil {
		interval -= 5
	}

	if p.Equipment.Armor != nil {
		interval -= 5
	}

	if p.Equipment.Helmet != nil {
		interval -= 2
	}

	if p.Equipment.Shield != nil {
		interval -= 5
	}

	if p.Equipment.Weapon != nil {
		interval -= 2
	}

	return interval
}

// аналогично с  водой
func (p *Player) GetThirstInterval() int {
	interval := 50

	if p.Equipment.Bag != nil {
		interval -= 3
	}

	if p.Equipment.Armor != nil {
		interval -= 3
	}

	if p.Equipment.Helmet != nil {
		interval -= 1
	}

	if p.Equipment.Shield != nil {
		interval -= 3
	}

	if p.Equipment.Weapon != nil {
		interval -= 1
	}

	return interval
}
