package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/monster"
	"log"
	"strings"

	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RoomModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	Exits       string `gorm:"type:text"`
	Items       string `gorm:"type:text"`
	MonsterData string `gorm:"type:text"`
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  gorm.DeletedAt `gorm:"index"`

	PlayerOccupantID string `gorm:"size:36"`
}

func (m *RoomModel) ToEntity() (*Room, error) {
	//парсим выходы из JSON (то, что в БД)
	var exits map[string]string
	if m.Exits != "" {
		err := json.Unmarshal([]byte(m.Exits), &exits)
		if err != nil {
			return nil, errors.New("Не удалось распарсить выходы")
		}
	}
	//парсим предметы
	var items []*item.ItemStack
	if m.Items != "" {
		err := json.Unmarshal([]byte(m.Items), &items)
		if err != nil {
			return nil, errors.New("Не удалось распарсить предметы")
		}
	}
	room := &Room{
		ID:               m.ID,
		Name:             m.Name,
		Description:      m.Description,
		Exits:            exits,
		Items:            items,
		TownExits:        []TownExit{}, // пока пусто
		playerOccupantID: m.PlayerOccupantID,
	}

	monsterData, err := m.getMonster()
	if err != nil {
		log.Printf("Ошибка создания монстра: %v", err)
	}

	room.Monster = monsterData
	if room.ID == "dungeon_goblin" && room.Monster == nil {
		room.Monster = monster.NewGoblin(room.ID)
	}
	if room.ID == "global_town" {
		for exitName, roomID := range room.Exits {
			// Проверяем, что это выход к дороге игрока
			if strings.Contains(roomID, "road_") {
				// Извлекаем OwnerID (часть после road_)
				parts := strings.Split(roomID, "_") // ["road", "player", "123"] Split: Разрезает строку на части по разделителю "_".
				if len(parts) >= 2 {
					// Собираем всё после первого "_"
					ownerID := strings.Join(parts[1:], "_") // "player_123" parts[1:] означает "все элементы, начиная с индекса 1" (пропускаем "road").

					room.TownExits = append(room.TownExits, TownExit{
						Name:    exitName,
						RoomID:  roomID,
						OwnerID: ownerID,
					})
				}
			}

		}
	}
	return room, nil
}

// вернуть комнату в бд
func FromEntity(r *Room) (*RoomModel, error) {
	//парс в JSON выходы
	exits, err := json.Marshal(r.Exits)
	if err != nil {
		return nil, errors.New("Не удалось преобразовать выходы")
	}

	//парс в JSON предметы из комнаты
	items, err := json.Marshal(r.Items)
	if err != nil {
		return nil, errors.New("Не удалось преобразовать предметы")
	}

	//создаем модель
	model := &RoomModel{
		ID:               r.ID,
		Name:             r.Name,
		Description:      r.Description,
		Exits:            string(exits),
		Items:            string(items),
		Created_at:       time.Now(),
		Updated_at:       time.Now(),
		PlayerOccupantID: r.playerOccupantID,
	}

	//сохраняем монстра
	if err := model.saveMonsterJSON(r.Monster); err != nil {
		return nil, err
	}
	return model, nil
}

// saveMonsterJSOM сохраняет монстра в JSON
func (m *RoomModel) saveMonsterJSON(mon *monster.Monster) error {

	if mon == nil {
		m.MonsterData = ""
		return nil
	}

	data, err := json.Marshal(mon)
	if err != nil {
		return err
	}

	m.MonsterData = string(data)
	return nil
}

// getMonster загрузка из JSON
func (m *RoomModel) getMonster() (*monster.Monster, error) {
	if m.MonsterData == "" {
		return nil, nil
	}

	var mon monster.Monster
	err := json.Unmarshal([]byte(m.MonsterData), &mon)
	return &mon, err
}
