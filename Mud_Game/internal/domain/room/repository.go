package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/monster"
)

// Room = сама комната (что она делает)
type RoomInterface interface { //интерфейс комнаты
	//минимальный набор(обязательный)
	GetID() string               //Возвращает уникальный ID комнаты	Чтобы знать, где игрок находится
	GetName() string             // Возвращает название ("Таверна", "Лес")	Для вывода игроку
	GetDescription() string      //Возвращает описание	Для команды look
	GetExits() map[string]string //Возвращает карту выходов Чтобы знать, куда можно идти

	//действия
	TakeItem(itemName string, count int) (*item.ItemStack, error) //Забирает предмет из комнаты (удаляет или уменьшает)
	GetItems() []*item.ItemStack                                  //Просто показывает, что лежит в комнате
	OnEnter(playerID string) string                               //Что происходит при входе игрока	Приветствие, события, ловушки
	OnExit(playerID string) string                                //Что происходит при выходе	Попрощаться, закрыть дверь
	Look(playerID string) string                                  //Осмотр комнаты	Описание + предметы + игроки
	AddItem(stack *item.ItemStack) error                          //скинуть предмет в комнату

	GetMonster() *monster.Monster  //наличие понстра в комнате
	SetMonster(m *monster.Monster) //обновление монстра(после урона)
}

// Repository = где комнаты лежат (как их найти/сохранить)
type Repository interface { //что умеет хранилище комнат
	FindByID(id string) (RoomInterface, error) //возвращает интерфейс RoomInterface, Ищет комнату по ID	Чтобы переместить игрока
	Save(room RoomInterface) error             //сохраняем любую комнату,Чтобы добавить новую или обновить
	Delete(id string) error                    //Удаляет комнату	Если комната больше не нужна
	FindAll() ([]RoomInterface, error)
}
