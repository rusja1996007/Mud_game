package memoryrepo

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"errors"
	"sync"
)

// in-memory игроки
type MemoryRepository struct {
	players map[string]*player.Player //Хранилище данных - клю это Будет ID, значение указатель на player.Player (* означает, что храним ссылку, а не копию)
	mtx     sync.RWMutex              //замок для безопасного доступа из разных горутин
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{ //создает и возвращает адрес созданой структуры
		players: make(map[string]*player.Player),
	}
}
func (x *MemoryRepository) Save(p *player.Player) error {
	//(x *MemoryRepository) это получатель (receiver). Метод привязан к конкретному экземпляру MemoryRepository

	if p == nil {
		return errors.New("Игрока нет")
	}
	//проверка в Save — это защита от дурака (от ошибок в коде сервера), а не требование к игроку.
	if p.ID == "" {
		return errors.New("ID игрока не может быть пустым")
	}
	if p.Name == "" {
		return errors.New("Введите имя игрока ")
	}

	x.mtx.Lock()
	defer x.mtx.Unlock()
	//_- Переменная, но она не жуна тут. Нужно узнать тру или фолс.
	if _, ok := x.players[p.ID]; ok {
		return errors.New("Игрок с таким ID существует")
	}
	x.players[p.ID] = p

	return nil

}
func (x *MemoryRepository) FindByID(id string) (*player.Player, error) {
	if id == "" {
		return nil, errors.New("ID не может быть пустым")
	}
	x.mtx.RLock()
	defer x.mtx.RUnlock()
	// player — переменная, кого искали (если нашли),ok — это флаг успеха (нашли/не нашли)
	//чтобы искать - Прямой доступ по ключу (x.players[id])
	player, ok := x.players[id]
	if !ok {
		return nil, errors.New("Игрок с таким ID не найден")
	}
	return player, nil

}
func (x *MemoryRepository) FindByName(name string) (*player.Player, error) {
	if name == "" {
		return nil, errors.New("Введите имя игрока")
	}
	x.mtx.RLock()
	defer x.mtx.RUnlock()
	//чтобы искать по name - Перебор всех (for ... range)
	for _, p := range x.players {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, errors.New("Игрок не найден")
}
func (x *MemoryRepository) Delete(id string) error {
	if id == "" {
		return errors.New("Введите ID")
	}

	x.mtx.Lock()
	defer x.mtx.Unlock()

	if _, ok := x.players[id]; !ok {
		return errors.New("Игрок с таким ID не найден")
	}

	delete(x.players, id)
	return nil

}
