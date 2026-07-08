package player

import (
	"Mud_game/Mud_Game/internal/domain/buff"
	"Mud_game/Mud_Game/internal/domain/combat"
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/loot"

	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// // Equipment - экипировка (не занимает слоты рюкзака)
type Equipment struct {
	Weapon   *item.ItemStack
	Armor    *item.ItemStack
	Helmet   *item.ItemStack
	Bag      *item.ItemStack   //мешок(сумка)
	Shield   *item.ItemStack   //щит
	Boots    *item.ItemStack   //тапки
	Ring1    *item.ItemStack   //кольцо1
	Ring2    *item.ItemStack   //кольцо2
	BagItems []*item.ItemStack //мешок как отдельный контейнер

}

// Управление (каналы, мьютексы, временные флаги)
type Player struct {
	connMutex sync.Mutex //для безопасной записи в conn

	ID          string
	Name        string
	CurrentRoom string            //текущая комната
	Inventory   []*item.ItemStack // ← стопки предметов (название + кол-во)
	Equipment   *Equipment
	Zone        *PLayerZone
	Stats       *Stats // характеристики

	//охота
	PendingHunt       bool      // ожидание подтверждения охото(yes)
	PendingHuntExpiry time.Time //время, когда запрос истекает
	stopHunger        chan bool //добавили чтобы не увеличивались тики после охоты(остановка тиков)
	stopThirst        chan bool

	// характеристики
	PendingStatChoiсe       bool      //ожидает ли игрок выбор
	PendingStatChoiсeExpiry time.Time //время на выбор характеристики

	//бафы
	ActiveBuffs    []*buff.Buff //список активных бафов
	stopBuffTicker chan bool    //канал остановки тикера

	//путешествие
	PendingTravel          bool
	PendingTravelDirection string    //Ожидаемое направление поездки
	PendingTravelExpiry    time.Time //Ожидающее истечение срока действия поездки

	//поиск
	IsSearching bool //сейчас обыскивает?(лут)

	//данж
	stopDungeonTimer chan bool //остановка таймера

}

// Данные (состояния, характеристики)
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
	Level            int
	Experience       int
	PendingStatPoint int //неиспользованые очки характеристик

	//охота
	IsHunting      bool      // на охоте ли персонаж
	HuntingEndTime time.Time //когда закончится охота

	//сон
	IsSleeping      bool      //спит дома
	SleepStartTime  time.Time //время когда начал спать
	IsSleepingHotel bool      //спит в отеле

	//путеществие
	IsTraveling      bool      //в путешествии?
	TravelEndTime    time.Time //время когда закончится путешествие
	TravelTargetRoom string    //Это поле хранит ID комнаты, куда игрок идёт:

	//временные бонусы :
	MaxHealthBonus int //к максимальному здоровью

	//данжи
	IsInDungeon      bool      //в данже?
	EnteredDungeonAt time.Time //вошел в подземелье в ...
}

// для сохранения игрока при отключении сервера без проблем с циклическим импортом
type Saver interface {
	Save(p *Player) error
}

// SendMessage безопасно отправляет сообщение игроку
func (p *Player) SendMessage(conn net.Conn, msg string) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()
	fmt.Fprint(conn, msg)
}

// запускает таймер голода (каждые GetHUNGERInterval секунд )
func (p *Player) StartHungerTicker(conn net.Conn, repo Repository) {
	if p.Stats.Health <= 0 {
		return
	}
	repo.Save(p)
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
				if p.Stats.Health > 0 {
					repo.Save(p) // сохраняем только живых
				}
			}
		}
	}()
}

// таймер жажды (каждые getthirstticker секунд -1)
func (p *Player) StartThirstTicker(conn net.Conn, repo Repository) {
	if p.Stats.Health <= 0 {
		return
	}
	repo.Save(p)
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
				if p.Stats.Health > 0 {
					repo.Save(p) // сохраняем только живых
				}
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
	defence := p.GetTotalDefence()
	huntLoot, wolfResult, brokenMsg, totalDamage, totalXP := loot.GenerateHuntLoot(weapon, p.Stats.Tracking, defence)

	//применяем урон
	if totalDamage > 0 {
		p.TakeDamage(totalDamage, combat.DamagePhysical, conn)
		if p.Stats.Health <= 0 {
			fmt.Fprintf(conn, "Ты погиб на охоте!\n")
			repo.Delete(p.ID)
			conn.Close()
			return
		}
	}

	if totalXP > 0 {
		p.AddExperience(totalXP, conn)
		fmt.Fprintf(conn, "Ты получил %d опыта.\n", totalXP)
	}
	if wolfResult != nil && wolfResult.Win && weapon != nil {
		if weapon.Decrease(5) {
			brokenMsg = "Твое оружие сломалось в бою!"
			p.Equipment.Weapon = nil
		}
	}

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

	p.CheckLevelUp(conn)
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

// возвращает общую физическую защиту игрока
func (p *Player) GetTotalDefence() int {
	total := 0

	if p.Equipment.Armor != nil && p.Equipment.Armor.Durability > 0 {
		total += p.Equipment.Armor.Defence
	}

	if p.Equipment.Bag != nil && p.Equipment.Bag.Durability > 0 {
		total += p.Equipment.Bag.Defence
	}

	if p.Equipment.Boots != nil && p.Equipment.Boots.Durability > 0 {
		total += p.Equipment.Boots.Defence
	}

	if p.Equipment.Shield != nil && p.Equipment.Shield.Durability > 0 {
		total += p.Equipment.Shield.Defence
	}
	if p.Equipment.Helmet != nil && p.Equipment.Helmet.Durability > 0 {
		total += p.Equipment.Helmet.Defence
	}

	return total

}

// возвращает общую огненую защиту игрока
func (p *Player) GetTotalFireDefence() int {
	total := 0

	if p.Equipment.Helmet != nil && p.Equipment.Helmet.Durability > 0 {
		total += p.Equipment.Helmet.FireDefence
	}

	if p.Equipment.Armor != nil && p.Equipment.Armor.Durability > 0 {
		total += p.Equipment.Armor.FireDefence
	}

	if p.Equipment.Bag != nil && p.Equipment.Bag.Durability > 0 {
		total += p.Equipment.Bag.FireDefence
	}

	if p.Equipment.Shield != nil && p.Equipment.Shield.Durability > 0 {
		total += p.Equipment.Shield.FireDefence
	}

	if p.Equipment.Boots != nil && p.Equipment.Boots.Durability > 0 {
		total += p.Equipment.Boots.FireDefence
	}

	if p.Equipment.Ring1 != nil {
		total += p.Equipment.Ring1.FireDefence
	}

	if p.Equipment.Ring2 != nil {
		total += p.Equipment.Ring2.FireDefence
	}
	return total
}

// возвращают общую защиту от яда
func (p *Player) GetTotalPoisonDefence() int {
	total := 0

	if p.Equipment.Helmet != nil && p.Equipment.Helmet.Durability > 0 {
		total += p.Equipment.Helmet.PoisonDefence
	}

	if p.Equipment.Armor != nil && p.Equipment.Armor.Durability > 0 {
		total += p.Equipment.Armor.PoisonDefence
	}

	if p.Equipment.Bag != nil && p.Equipment.Bag.Durability > 0 {
		total += p.Equipment.Bag.PoisonDefence
	}

	if p.Equipment.Shield != nil && p.Equipment.Shield.Durability > 0 {
		total += p.Equipment.Shield.PoisonDefence
	}

	if p.Equipment.Boots != nil && p.Equipment.Boots.Durability > 0 {
		total += p.Equipment.Boots.PoisonDefence
	}

	if p.Equipment.Ring1 != nil {
		total += p.Equipment.Ring1.PoisonDefence
	}

	if p.Equipment.Ring2 != nil {
		total += p.Equipment.Ring2.PoisonDefence
	}
	return total
}

func (p *Player) GetTotalMagicDefence() int {
	total := 0

	if p.Equipment.Helmet != nil && p.Equipment.Helmet.Durability > 0 {
		total += p.Equipment.Helmet.MagicDefence
	}

	if p.Equipment.Armor != nil && p.Equipment.Armor.Durability > 0 {
		total += p.Equipment.Armor.MagicDefence
	}

	if p.Equipment.Bag != nil && p.Equipment.Bag.Durability > 0 {
		total += p.Equipment.Bag.MagicDefence
	}

	if p.Equipment.Shield != nil && p.Equipment.Shield.Durability > 0 {
		total += p.Equipment.Shield.MagicDefence
	}

	if p.Equipment.Boots != nil && p.Equipment.Boots.Durability > 0 {
		total += p.Equipment.Boots.MagicDefence
	}

	if p.Equipment.Ring1 != nil {
		total += p.Equipment.Ring1.MagicDefence
	}

	if p.Equipment.Ring2 != nil {
		total += p.Equipment.Ring2.MagicDefence
	}
	return total
}

// наносит урон игроку с учетом защиты
func (p *Player) TakeDamage(damage int, dmgType combat.DamageType, conn net.Conn) {
	defence := 0

	switch dmgType {
	case combat.DamagePhysical:
		defence = p.GetTotalDefence()
	case combat.DamageFire:
		defence = p.GetTotalFireDefence()
	case combat.DamagePoison:
		defence = p.GetTotalPoisonDefence()
	case combat.DamageMagic:
		defence = p.GetTotalMagicDefence()
	default:
		defence = 0
	}

	//процентное снижение урона (защита/защита+100)
	reduction := float64(defence) / (float64(defence) + 100)
	finalDamage := int(float64(damage) * (1 - reduction))
	finalDamage = max(finalDamage, 1) //возвращает большее из двух чисел(замена если fdamage<1-fdamage=1)

	p.Stats.Health -= finalDamage

	fmt.Fprintf(conn, "Ты получил %d урона!\n", finalDamage)

	p.DecreaseArmorDurability()
}

// снижаем прочность брони, 20% шанс износа при каждом ударе
func (p *Player) DecreaseArmorDurability() {
	if rand.Intn(100) >= 50 { ///////////////////////////////////
		return //не изнашивается
	}

	//износ...
	if p.Equipment.Helmet != nil {
		if p.Equipment.Helmet.Decrease(2) {
			fmt.Println("Твой шлем сломался!")
			p.Equipment.Helmet = nil
		}
	}

	if p.Equipment.Armor != nil {
		if p.Equipment.Armor.Decrease(4) {
			fmt.Println("Твоя броня сломалась!")
			p.Equipment.Armor = nil
		}
	}

	if p.Equipment.Shield != nil {
		if p.Equipment.Shield.Decrease(8) {
			fmt.Println("Твой щит сломался!")
			p.Equipment.Shield = nil
		}
	}

	if p.Equipment.Boots != nil {
		if p.Equipment.Boots.Decrease(2) {
			fmt.Println("Твои ботинки уничтожены!")
			p.Equipment.Boots = nil
		}
	}

}

// сколько опыта нужно для следующего уровня.
func (p *Player) GetExpForNextLevel() int {
	return 100 * p.Stats.Level
}

func (p *Player) AddExperience(amount int, conn net.Conn) {
	p.Stats.Experience += amount
	p.CheckLevelUp(conn)
}

func (p *Player) LevelUp(conn net.Conn) {
	p.Stats.PendingStatPoint++
	p.Stats.Level++
	fmt.Fprintf(conn, "\n🎉 Поздравляем! Вы достигли %d уровня!\n", p.Stats.Level)
	fmt.Fprintf(conn, "Вы получили очко характеристик\n")
	fmt.Fprintf(conn, "Используйте команду 'statpoints' для распределения\n")

}

func (p *Player) CheckLevelUp(conn net.Conn) bool {
	leveled := false
	for p.Stats.Experience >= p.GetExpForNextLevel() {
		p.Stats.Experience -= p.GetExpForNextLevel()
		p.LevelUp(conn)
		leveled = true

	}
	return leveled
}

// процес выбора характеристики при лвлАпе
func (p *Player) ProcessStatChoice(choice string, conn net.Conn) {
	if !p.PendingStatChoiсe {
		return
	}

	p.PendingStatChoiсe = false
	p.Stats.PendingStatPoint--

	switch choice {
	case "1":
		p.Stats.Strength++
		p.Stats.Health = 50 + p.Stats.Strength*5
		fmt.Fprintf(conn, "Сила увеличена до %d, максимальное здоровье: %d\n", p.Stats.Strength, p.Stats.Health)
	case "2":
		p.Stats.Dexterity++
		fmt.Fprintf(conn, "Ловкость увеличена до %d\n", p.Stats.Dexterity)
	case "3":
		p.Stats.Intelect++
		fmt.Fprintf(conn, "Интеллект увеличен до %d\n", p.Stats.Intelect)
	case "4":
		p.Stats.Tracking++
		fmt.Fprintf(conn, "Следопытство увеличено до %d\n", p.Stats.Tracking)
	}
	fmt.Fprintf(conn, " >")

}

// запуск бафа
func (p *Player) StartBuffTicker(conn net.Conn, repo Repository) {
	if p.Stats.Health <= 0 {
		return
	}

	if p.stopBuffTicker != nil {
		select {
		case p.stopBuffTicker <- true:
		default:
			close(p.stopBuffTicker)

		}
	}

	p.stopBuffTicker = make(chan bool)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-p.stopBuffTicker:
				return
			case <-ticker.C:
				for _, buff := range p.ActiveBuffs {
					if buff.RemainingTime > 0 {
						buff.RemainingTime -= 1 * time.Second
					}
				}
				p.processBuffs(repo)

				if p.Stats.Health > 0 {
					repo.Save(p)
				}
			}
		}
	}()
}

// Проверяет все активные баффы игрока и применяет их эффекты.
func (p *Player) processBuffs(repo Repository) {
	now := time.Now()
	newBuffs := make([]*buff.Buff, 0)

	for _, b := range p.ActiveBuffs {
		if b.RemainingTime <= 0 {

			if b.Type == buff.MaxHealthBoost {
				p.Stats.MaxHealthBonus -= b.Value
				maxHealt := 50 + p.Stats.Strength*5 + p.Stats.MaxHealthBonus
				if p.Stats.Health > maxHealt {
					p.Stats.Health = maxHealt
				}
			}
			continue // бафф истёк
		}

		// Проверяем интервал срабатывания
		if b.Interval > 0 {
			// Если LastTick не установлен или прошло >= Interval
			if b.LastTick.IsZero() || now.Sub(b.LastTick) >= b.Interval {
				p.ApplyBuffEffect(b)
				b.LastTick = now
			}
		}

		newBuffs = append(newBuffs, b)
	}

	p.ActiveBuffs = newBuffs
	repo.Save(p)
}

// применить баф
func (p *Player) ApplyBuffEffect(b *buff.Buff) {
	switch b.Type {
	case buff.HealthRegen:
		p.Stats.Health += b.Value
		maxHealth := 50 + p.Stats.Strength*5
		if p.Stats.Health > maxHealth {
			p.Stats.Health = maxHealth
		}
		//fmt.Fprintf(conn, "Вы восстановили %d жизней.\n", b.Value)

	case buff.MaxHealthBoost:
		p.Stats.MaxHealthBonus += b.Value
		maxHealt := 50 + p.Stats.Strength*5 + p.Stats.MaxHealthBonus
		if p.Stats.Health > maxHealt {
			p.Stats.Health = maxHealt
		}

	}
}

// StopAllTickers останавливает все тикеры игрока
func (p *Player) StopAllTickers() {
	if p.stopBuffTicker != nil {
		select {
		case p.stopBuffTicker <- true:
		default:
		}
		close(p.stopBuffTicker)
		p.stopBuffTicker = nil
	}
	if p.stopHunger != nil {
		select {
		case p.stopHunger <- true:
		default:
		}
		close(p.stopHunger)
		p.stopHunger = nil
	}
	if p.stopThirst != nil {
		select {
		case p.stopThirst <- true:
		default:
		}
		close(p.stopThirst)
		p.stopThirst = nil
	}
}

// кик игрока при долгом присутствии в данже
func (p *Player) StartDungeonKickTimer(conn net.Conn, repo Repository, roomRepo room.Repository) {

	if p.Stats.IsInDungeon {

		p.stopDungeonTimer = make(chan bool)

		select {
		case <-p.stopDungeonTimer:
			return //остановили
		case <-time.After(1 * time.Minute): ///////////////////////для теста через это время сработает этот кик:

			//текущая комната
			room, _ := roomRepo.FindByID(p.CurrentRoom)
			monster := room.GetMonster()

			//если монстр жив-автоматический побег с получением урона
			if monster != nil && monster.IsAlive {
				monsterDamage := monster.MinDamage + rand.Intn(monster.MaxDamage-monster.MinDamage+1)
				defence := p.GetTotalDefence()
				reduction := float64(defence) / (float64(defence) + 100)
				finalDamage := int(float64(monsterDamage) * (1 - reduction))
				if finalDamage <= 0 {
					finalDamage = 1
				}

				p.Stats.Health -= finalDamage
				monster.Health = monster.MaxHealth
				if p.Stats.Health > 0 {
					repo.Save(p)
				}
				msg := fmt.Sprintf("💨 Ты слишком долго был в бою и сбежал! Монстр нанёс %d урона вслед.\n> ", finalDamage)
				p.SendMessage(conn, msg)

				if p.Stats.Health <= 0 {
					p.SendMessage(conn, "💀Ты погиб...\n")
					repo.Delete(p.ID)
					conn.Close()
					return
				}

				// телепорт
				p.CurrentRoom = room.GetExitRoomID()
				room.SetPlayerOccupantID("")
				p.Stats.IsInDungeon = false
				p.Stats.EnteredDungeonAt = time.Time{}
				roomRepo.Save(room)
				if p.Stats.Health > 0 {
					repo.Save(p)
				}
				return
			}
			// Если монстр мёртв — телепорт
			if room.GetExitRoomID() != "" {
				p.CurrentRoom = room.GetExitRoomID()
				room.SetPlayerOccupantID("")
				p.Stats.IsInDungeon = false
				p.Stats.EnteredDungeonAt = time.Time{}

				roomRepo.Save(room)
				if p.Stats.Health > 0 {
					repo.Save(p)
				}

				p.SendMessage(conn, "⏰Вы слишком долго были в подземелье, вас выкинуло.\n> ")

			}
		}

	}
}

// StopDungeonTimer останавливает таймер кика из данжа
func (p *Player) StopDungeonTimer() {
	if p.stopDungeonTimer != nil {
		p.stopDungeonTimer <- true
		close(p.stopDungeonTimer)
		p.stopDungeonTimer = nil
	}
}

// побег из подземелья если "лег" сервер(безопасный)
func (p *Player) HandleDisconnect(saver Saver, roomRepo room.Repository) {

	if p.Stats.Health <= 0 {
		return // ← должно быть, чтобы не сохранять мёртвого
	}
	if p.CurrentRoom == "dungeon_goblin" {
		room, _ := roomRepo.FindByID(p.CurrentRoom)
		//телепорт
		p.CurrentRoom = room.GetExitRoomID()
		room.SetPlayerOccupantID("")
		p.Stats.IsInDungeon = false
		p.Stats.EnteredDungeonAt = time.Time{}
		roomRepo.Save(room)
		saver.Save(p)

	}
}

// шанс на поломку(при выживании)
func (p *Player) BreakAllEquipment() {
	slots := []**item.ItemStack{
		&p.Equipment.Weapon,
		&p.Equipment.Armor,
		&p.Equipment.Helmet,
		&p.Equipment.Shield,
		&p.Equipment.Boots,
		&p.Equipment.Bag,
		&p.Equipment.Ring1,
		&p.Equipment.Ring2,
	}
	for _, slot := range slots {
		if rand.Intn(100) < 66 {
			if *slot != nil {
				(*slot).Durability = 0
			}
		}
	}
}
