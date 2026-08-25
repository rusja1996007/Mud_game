package trader

import (
	"Mud_game/Mud_Game/internal/domain/npc"
	"errors"

	"gorm.io/gorm"
)

type PostgresTraderRepository struct {
	db *gorm.DB
}

// конструктор
func NewPostgresTraderRepository(db *gorm.DB) (*PostgresTraderRepository, error) {
	if err := db.AutoMigrate(&npc.TraderModel{}); err != nil { //создание таблиц с traderModel
		return nil, err
	}

	return &PostgresTraderRepository{
		db: db,
	}, nil
}

// сохранить торговца
func (p *PostgresTraderRepository) Save(t *npc.Trader) error {
	model, err := npc.FromEntity(t)
	if err != nil {
		return err
	}

	result := p.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// найти по ID торговца
func (p *PostgresTraderRepository) FindByID(id string) (*npc.Trader, error) {
	model := npc.TraderModel{}

	if err := p.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	//конвертир модель в сущность
	traderEntity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}
	return traderEntity, nil
}

// найти  в комнате торговца
func (p *PostgresTraderRepository) FindByRoom(roomId string) ([]*npc.Trader, error) {
	model := []*npc.TraderModel{} //слайс моделей

	err := p.db.Where("room_id = ?", roomId).Find(&model).Error
	if err != nil {
		return nil, err
	}

	//превращаем модели в сущности все
	var traders []*npc.Trader

	for _, model := range model {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		traders = append(traders, entity)

	}

	return traders, nil
}

// удалить торговца (пока ненадо но пусть будет)
func (p *PostgresTraderRepository) Delete(id string) error {

	result := p.db.Delete(&npc.TraderModel{}, "id = ?", id)

	return result.Error
}
