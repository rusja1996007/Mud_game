package room

import (
	"Mud_game/Mud_Game/internal/domain/item"
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
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  gorm.DeletedAt `gorm:"index"`
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
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Exits:       exits,
		Items:       items,
		TownExits:   []TownExit{}, // пока пусто
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
	return &RoomModel{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Exits:       string(exits),
		Items:       string(items),
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}, nil
}
