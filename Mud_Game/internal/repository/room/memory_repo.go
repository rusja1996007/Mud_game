package room

import (
	"Mud_game/Mud_Game/internal/domain/room"
	"errors"
	"sync"
)

type MemoryRepository struct {
	rooms map[string]room.RoomInterface //храним интерфейс
	mtx   sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{ //создает и возвращает адрес созданой структуры
		rooms: make(map[string]room.RoomInterface),
	}
}
func (x *MemoryRepository) FindByID(id string) (room.RoomInterface, error) {
	if id == "" {
		return nil, errors.New("Введите ID комнаты")
	}
	x.mtx.RLock()
	defer x.mtx.RUnlock()

	rm, ok := x.rooms[id]
	if !ok {
		return nil, errors.New("Комната по данному ID не существует")
	}
	return rm, nil

}
func (x *MemoryRepository) Save(r room.RoomInterface) error { //сохранить комнату
	if r == nil {
		return errors.New("Введите интерфейс комнаты")
	}

	id := r.GetID() //сделали чтобы проверить что id не пустой
	if id == "" {
		return errors.New("Поле ID пустое")
	}

	x.mtx.Lock()
	defer x.mtx.Unlock()

	x.rooms[r.GetID()] = r //сохраняем мапу в MemoryRepository
	return nil

}
func (x *MemoryRepository) Delete(id string) error { //удалить комнату
	if id == "" {
		return errors.New("Введите ID комнаты")
	}
	x.mtx.Lock()
	defer x.mtx.Unlock()

	_, ok := x.rooms[id]
	if !ok {
		return errors.New("ID не найдено")
	}
	delete(x.rooms, id)
	return nil
}
func (r *MemoryRepository) FindAll() ([]room.RoomInterface, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	rooms := make([]room.RoomInterface, 0, len(r.rooms))
	for _, rm := range r.rooms {
		rooms = append(rooms, rm)
	}
	return rooms, nil
}
