package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/monster"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// TownExit - информация о выходе из города
type TownExit struct {
	Name    string //"дом Иван"- название выхода
	RoomID  string // "road_player_123"- ID дороги игрока
	OwnerID string // "player_123"- КОМУ принадлежит этот выход
}

type Room struct {
	ID          string
	Name        string
	Description string
	Exits       map[string]string //выходы: направление → ID комнаты
	Items       []*item.ItemStack //[]item.ItemStack = много разных предметов с количеством
	mtx         sync.RWMutex
	TownExits   []TownExit         `json:"-"` //Этот тег говорит GORM не сохранять это поле в БД.
	Monster     *monster.Monster   `json:"-"`
	MonsterS    []*monster.Monster `json:"-"` //НЕСКОЛЬКО

	playerOccupantID string    `json:"-"` //для блокировки данжа если там ктото есть
	NextSpawnTime    time.Time //следующее время обновления предметов

	ExitRoom string `json:"exit_room_id"` //куда телепортироваться при побеге
}

func (r *Room) GetID() string {
	return r.ID
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) GetDescription() string {
	return r.Description
}

func (r *Room) GetExits() map[string]string {
	return r.Exits
}
func (r *Room) Look(playerID string) string {
	//result += text	Каждый раз создаётся новая строка, старая копируется и удаляется.
	//Это нагружает память.
	var builder strings.Builder //Выделяем место для сборки строки.

	builder.WriteString(r.Name) //Текст добавляется в буфер, без копирования всей строки заново
	builder.WriteString("\n")

	//динамическое описание монстра
	aliveMonsters := r.GetAliveMonsters()
	if len(aliveMonsters) > 0 {
		for i, m := range aliveMonsters {
			builder.WriteString(fmt.Sprintf(" %d. %s (%d/%d)\n", i+1, m.Description, m.Health, m.MaxHealth))

		}
	} else if len(aliveMonsters) == 0 && (r.ID == "dungeon_goblin" || r.ID == "dungeon_goblins_v2") {
		// Проверяем таймер обвала, если монстры мертвы
		monster := r.GetMonster()
		if monster != nil && !monster.IsAlive {
			remaining := time.Until(monster.TimeToLoot).Round(time.Second)
			if remaining > 0 {
				builder.WriteString(fmt.Sprintf("⚠️ Пещера обвалится через %v. У тебя есть время на обыск!\n", remaining))
			} else {
				builder.WriteString("Пещера обвалилась.\n")
			}
		} else {
			builder.WriteString("Пещера пуста. Следы битвы видны повсюду.\n")
		}
	}

	items := r.GetItems()
	if len(items) > 0 {
		builder.WriteString("Вы видите:\n")
		for i, stack := range items {
			fmt.Fprintf(&builder, "  %d. %s", i+1, stack.Name)
			if stack.Durability > 0 {
				fmt.Fprintf(&builder, "(прочность:%d)", stack.Durability)
			}
			if stack.Count > 1 {
				fmt.Fprintf(&builder, " x%d", stack.Count)
			}

			builder.WriteString("\n")
		}
	} else {
		builder.WriteString("Вы не видите ничего интересного.\n")
	}

	builder.WriteString("Выходы: ")

	if r.ID == "global_town" { //// 1. Проверяем, что это город

		//показываем обычные выходы:
		for exitName := range r.Exits {
			if !strings.HasPrefix(exitName, "дом ") {
				builder.WriteString(exitName + ", ")
			}
		}

		//показываем персональные выходы(дом игрока)

		for _, exit := range r.TownExits { // 2. Проходим по ВСЕМ выходам города
			if exit.OwnerID == playerID { // 3. Сравниваем владельца с текущим игроком
				builder.WriteString(exit.Name + " ")
			}
		}
	} else {
		for exits := range r.Exits {
			builder.WriteString(exits)
			builder.WriteString(" ")
		}
	}
	return builder.String() //Строитель отдаёт всё, что мы насобирали, одной строкой.
}

func (r *Room) OnEnter(playerID string) string {
	return "Ты вошел в " + r.Name
}
func (r *Room) OnExit(playerID string) string {
	return "Ты покинул " + r.Name
}
func (r *Room) GetItems() []*item.ItemStack {
	return r.Items
}
func (r *Room) TakeItem(itemName string, count int) (*item.ItemStack, error) {
	// Сначала ищем предмет
	foundIndex := -1
	for i, stack := range r.Items {
		if stack.Name == itemName {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return nil, errors.New("Предмет не найден")
	}

	// Теперь работаем с найденным
	if r.Items[foundIndex].Count < count {
		return nil, errors.New("Недостаточно предметов")
	}

	// // Сохраняем информацию о предмете ДО изменения
	originalItem := r.Items[foundIndex]

	//Уменьшаем количество или удаляем
	r.Items[foundIndex].Count -= count
	//Удаление элемента из слайса
	//У тебя есть ряд печенек. Ты хочешь убрать одну (с индексом 1):
	//Берёшь все, кто до неё
	//Берёшь все, кто после неё
	//Складываешь вместе — и вуаля, средняя исчезла!
	if r.Items[foundIndex].Count == 0 {
		r.Items = append(r.Items[:foundIndex], r.Items[foundIndex+1:]...)
	}
	//возвращаем стопку с тем что взяли
	return &item.ItemStack{
		Name:          itemName,
		Count:         count,
		ItemType:      originalItem.ItemType,
		SlotBonus:     originalItem.SlotBonus,
		HungerRestore: originalItem.HungerRestore,
		ThirstRestore: originalItem.ThirstRestore,
		MinDamage:     originalItem.MinDamage,
		MaxDamage:     originalItem.MaxDamage,
		Durability:    originalItem.Durability,
		Description:   originalItem.Description,
		ID:            originalItem.ID,
		Defence:       originalItem.Defence,
		MagicDefence:  originalItem.MagicDefence,
		PoisonDefence: originalItem.PoisonDefence,
		FireDefence:   originalItem.FireDefence,
		HealMin:       originalItem.HealMin,
		HealMax:       originalItem.HealMax,
	}, nil
}
func (r *Room) AddItem(stack *item.ItemStack) error {
	if stack == nil {
		return errors.New("Стопка предметов не может быть пустой")
	}
	if stack.Name == "" {
		return errors.New("Название предмета не может быть пустым ")
	}
	if stack.Count <= 0 {
		return errors.New("Количество должно быть положительным")
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	// Ищем существующую стопку с таким же названием
	for i := range r.Items {
		if r.Items[i].CanStackWith(stack) {
			r.Items[i].Count += stack.Count
			return nil
		}
	}
	// Если не нашли - добавляем новую стопку
	r.Items = append(r.Items, stack)
	return nil

}

// проверка наличия монстра(или монстров)
func (r *Room) GetMonster() *monster.Monster {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	// Сначала проверяем Monsters (новый данж)
	if len(r.MonsterS) > 0 {
		for _, m := range r.MonsterS {
			if m.IsAlive {
				if m.CheckRespawn() {
					delete(r.Exits, "up")
				}
				return m
			}
		}
	}

	// Потом проверяем Monster (старый данж)
	if r.Monster != nil {
		if r.Monster.CheckRespawn() {
			delete(r.Exits, "up")
		}
		return r.Monster
	}

	return nil
}

// обновление монстра(после урона)
func (r *Room) SetMonster(m *monster.Monster) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Monster = m
}

// получение ID игрока который в комнате
func (r *Room) GetPlayerOccupantID() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.playerOccupantID
}

// установка ID игрока в комнате
func (r *Room) SetPlayerOccupantID(id string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.playerOccupantID = id
}

// очищение от предметов комнаты
func (r *Room) ClearItems() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Items = []*item.ItemStack{}
}

// генерация для входа в подземелье итемов
func (r *Room) RegenerateItems() {
	r.Items = []*item.ItemStack{}

	possibleItems := []string{"clover", "burdock", "boletus edulis"}

	for _, itemName := range possibleItems {

		// 30% шанс появления предмета
		if rand.Intn(100) < 30 {
			count := 1 + rand.Intn(3)

			newItem := item.GetItem(itemName, count)
			if newItem != nil {
				r.Items = append(r.Items, newItem)
			}
		}
	}
}

// куда телепортироваться при побеге
func (r *Room) GetExitRoomID() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.ExitRoom
}

func (r *Room) GetAliveMonsters() []*monster.Monster {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var alive []*monster.Monster

	// Если есть поле Monsters (новый данж) — используем его
	if len(r.MonsterS) > 0 {
		for _, m := range r.MonsterS {
			if m.IsAlive {
				alive = append(alive, m)
			}
		}
		return alive
	}

	// Если есть поле Monster (старый данж) — используем его
	if r.Monster != nil && r.Monster.IsAlive {

		alive = append(alive, r.Monster)
	}

	return alive
}

// GetMonsters возвращает всех монстров
func (r *Room) GetMonsters() []*monster.Monster {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if r.Monster != nil {
		return []*monster.Monster{r.Monster}
	}
	return r.MonsterS
}
