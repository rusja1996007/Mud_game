package player

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/loot"
	"Mud_game/Mud_Game/internal/domain/room"
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
	PendingHunt bool   // ожидание подтверждения охото(yes)

	stopHunger chan bool //добавили чтобы не увеличивались тики после охоты(остановка тиков)
	stopThirst chan bool
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

	//охота
	IsHunting      bool      // на охоте ли персонаж
	HuntingEndTime time.Time //когда закончится охота
}

// запускает таймер голода (каждые GetHUNGERInterval секунд )
func (p *Player) StartHungerTicker(conn net.Conn, repo Repository) {

	//останавливаем тикеры если есть
	if p.stopHunger != nil {
		select {
		case p.stopHunger <- true: //Если канал существует (!= nil), отправляем сигнал true
		default:
		}
		close(p.stopHunger) //Закрываем канал (close)
	}

	p.stopHunger = make(chan bool)

	go func() {
		for {
			select {
			case <-p.stopHunger: //если пришел сигнал выходим из горутины
				return
			case <-time.After(time.Duration(p.GetHungerInterval()) * time.Second):

				//уменьшаем голод, только если не на охоте
				if p.Stats.IsHunting {
					continue
				}

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
	}()
}

// таймер жажды (каждые getthirstticker секунд -1)
func (p *Player) StartThirstTicker(conn net.Conn, repo Repository) {
	if p.stopThirst != nil {
		select {
		case p.stopThirst <- true:
		default:
		}

		close(p.stopThirst)
	}
	p.stopThirst = make(chan bool)

	go func() {
		for {
			select {
			case <-p.stopThirst:
				return
			case <-time.After(time.Duration(p.GetThirstInterval()) * time.Second):
				if p.Stats.IsHunting {
					continue
				}

				if p.Stats.Thirst > 0 {
					p.Stats.Thirst--
				}

				if p.Stats.Thirst == 0 {
					p.Stats.Health -= 2
					fmt.Fprintf(conn, "Ты умираешь от жажды!\n> ")
				}
				if p.Stats.Health <= 0 {
					fmt.Fprintf(conn, "Ты погиб от обезвоживания, персонаж удаляется!\n> ")
					repo.Delete(p.ID)
					conn.Close()
					return
				}
				repo.Save(p) // сохраняем только живых
			}
		}
	}()
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

// старт охоты
func (p *Player) StartHunt(conn net.Conn, repo Repository, roomRepo room.Repository) {
	//тратим 2 бутылки
	RemoveItem(&p.Inventory, "water bottle", 2)

	//устанавливаем состояние охоты
	p.Stats.IsHunting = true

	//////////ВРЕМЕННО ДЛЯ ТЕСТА
	huntDuration := 10 * time.Second
	p.Stats.HuntingEndTime = time.Now().Add(huntDuration) ///////////////!!!!

	repo.Save(p)

	//запускаем таймер окончания охоты
	go func() {
		time.Sleep(huntDuration) //////////////////////
		p.EndHunt(conn, repo, roomRepo)
	}()

}

// завершение охоты
func (p *Player) EndHunt(conn net.Conn, repo Repository, roomRepo room.Repository) {

	if !p.Stats.IsHunting {
		return
	}
	//останавливаем тикеры
	if p.stopHunger != nil {
		close(p.stopHunger)
		p.stopHunger = nil
	}

	if p.stopThirst != nil {
		close(p.stopThirst)
		p.stopThirst = nil
	}

	p.Stats.Hunger = 20
	p.Stats.Thirst = 20

	// Генерируем лут
	weapon := p.Equipment.Weapon
	huntLoot, wolfResult, brokenMsg := loot.GenerateHuntLoot(weapon, p.Stats.Tracking)
	if wolfResult != nil && wolfResult.Message != "" {
		fmt.Fprintf(conn, "%s\n", wolfResult.Message)
	}
	if brokenMsg != "" {
		fmt.Fprintf(conn, "%s\n", brokenMsg)
	}

	// ВЫВОДИМ ЛУТ ИГРОКУ
	fmt.Fprintf(conn, "\n Ты вернулся с охоты!\n")
	fmt.Fprintf(conn, "Голод: %d/100, Жажда: %d/100\n", p.Stats.Hunger, p.Stats.Thirst)

	if len(huntLoot) > 0 {
		fmt.Fprintf(conn, "Ты нашёл:\n")
		for _, l := range huntLoot {
			if l.Count > 1 {
				fmt.Fprintf(conn, "  • %s x%d\n", l.Name, l.Count)
			} else {
				fmt.Fprintf(conn, "  • %s\n", l.Name)
			}
		}
		fmt.Fprintf(conn, "Все занёс в дом.\n")

	} else {
		fmt.Fprintf(conn, "К сожалению, тебе ничего не попалось.\n")
	}

	p.AddItemToInventory(item.GetItem("empty bottle", 2))

	// Находим дом
	homeRoom, err := roomRepo.FindByID(p.Zone.HomeRoomID)
	if err == nil {
		for _, lootItem := range huntLoot {
			homeRoom.AddItem(lootItem)
		}
		roomRepo.Save(homeRoom)
	}

	p.Stats.IsHunting = false

	go p.StartHungerTicker(conn, repo)
	go p.StartThirstTicker(conn, repo)

	repo.Save(p)
	fmt.Fprintf(conn, "> ")
}

// уменьшение прочности на amount единиц
func (p *Player) DecreaseWeaponDurability(amount int) bool {
	if p.Equipment.Weapon != nil {
		if p.Equipment.Weapon.Decrease(amount) {
			p.Equipment.Weapon = nil //сломался
			return true
		}
	}
	return false
}
