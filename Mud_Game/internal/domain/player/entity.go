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

// запускает таймер голода (каждые 60 секунд -1)
func (p *Player) StartHungerTicker(conn net.Conn, repo Repository) {

	//будильник каждые 60 секунд звенит 1 раз  постоянно
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	//Каждый раз, когда будильник звенит(.C), делай то, что внутри фигурных скобок".
	for range ticker.C {

		// этот код выполняется каждые 60 секунд
		p.Stats.Hunger--
		if p.Stats.Hunger <= 0 {
			p.Stats.Hunger = 0
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

// таймер жажды (каждые 40 секунд -1)
func (p *Player) StartThirstTicker(conn net.Conn, repo Repository) {
	ticker := time.NewTicker(40 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		p.Stats.Thirst--
		if p.Stats.Thirst <= 0 {
			p.Stats.Thirst = 0
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
