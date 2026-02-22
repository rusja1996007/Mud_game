package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"strconv"
	"strings"
)

type Room struct {
	ID          string
	Name        string
	Description string
	Exits       map[string]string //выходы: направление → ID комнаты
	Items       []item.ItemStack  //[]item.ItemStack = много разных предметов с количеством
	//item.ItemStack = один вид предметов
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
