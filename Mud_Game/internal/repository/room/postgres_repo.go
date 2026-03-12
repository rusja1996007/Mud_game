package room

import (
	"Mud_game/Mud_Game/internal/domain/room"
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) (*PostgresRepository, error) {
	if err := db.AutoMigrate(&room.RoomModel{}); err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db: db,
	}, nil
}
func (pr *PostgresRepository) Save(r room.RoomInterface) error {

	// Нам нужно получить *Room из интерфейса
	// Но RoomInterface не имеет метода для получения конкретного типа

	// Придётся делать приведение типа:

	roomEntity, ok := r.(*room.Room)
	if !ok {
		return errors.New("Передан неверный тип комнаты")
	}

	model, err := room.FromEntity(roomEntity)
	if err != nil {
		return err
	}
	result := pr.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (pr *PostgresRepository) FindByID(id string) (room.RoomInterface, error) {
	model := room.RoomModel{}
	if err := pr.db.First(&model, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //комната не найдена
		}
		return nil, err //другая ошибка
	}
	return model.ToEntity()

}

// Этот метод загружает ВСЕ комнаты из БД при старте сервера
func (pr *PostgresRepository) FindAll() ([]room.RoomInterface, error) {
	var models []room.RoomModel
	if err := pr.db.Find(&models).Error; err != nil { //SELECT * FROM room_models;
		return nil, err
	}

	rooms := make([]room.RoomInterface, len(models))
	//Проходим по всем моделям из БД и превращаем каждую в настоящую комнату, с которой работает игра.
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		rooms[i] = entity
	}
	return rooms, nil

}
func (pr *PostgresRepository) Delete(id string) error {
	// Unscoped() = игнорируй мягкое удаление, удаляй насовсем
	return pr.db.Unscoped().Delete(&room.RoomModel{}, "id=?", id).Error
}
