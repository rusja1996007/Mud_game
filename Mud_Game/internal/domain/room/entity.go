package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Room struct {
	ID          string
	Name        string
	Description string
	Exits       map[string]string //выходы: направление → ID комнаты
	Items       []*item.ItemStack //[]item.ItemStack = много разных предметов с количеством
	mtx         sync.RWMutex
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
	builder.WriteString(r.Description)
	builder.WriteString("\n")

	if len(r.Items) > 0 {
		builder.WriteString("В доме есть: ")
		for i, stack := range r.Items {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString("(")
			builder.WriteString(strconv.Itoa(stack.Count))
			builder.WriteString(")")
			builder.WriteString(stack.Name)
		}
		builder.WriteString("\n")
	}
	builder.WriteString("Выходы: ")
	for exits := range r.Exits {
		builder.WriteString(exits)
		builder.WriteString(" ")
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
		Name:  itemName,
		Count: count,
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
		if r.Items[i].Name == stack.Name {
			r.Items[i].Count += stack.Count
			return nil
		}
	}
	// Если не нашли - добавляем новую стопку
	r.Items = append(r.Items, stack)
	return nil

}
