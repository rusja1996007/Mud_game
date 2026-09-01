package npc_repo

import (
	"Mud_game/Mud_Game/internal/domain/npc"
	"errors"

	"gorm.io/gorm"
)

type PostgresNPCRepository struct {
	db *gorm.DB
}

// конструктор
func NewPostgresNPCRepository(db *gorm.DB) (*PostgresNPCRepository, error) {
	if err := db.AutoMigrate(&npc.NPCModel{}); err != nil { //создание таблиц с traderModel
		return nil, err
	}

	repo := &PostgresNPCRepository{
		db: db,
	}

	//инициализируем NPC
	if err := repo.initNPCs(); err != nil {
		return nil, err
	}

	return repo, nil

}

// сохранить NPC
func (p *PostgresNPCRepository) Save(n *npc.NPC) error {
	model, err := npc.FromEntity(n)
	if err != nil {
		return err
	}

	result := p.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// найти по ID NPC
func (p *PostgresNPCRepository) FindByID(id string) (*npc.NPC, error) {
	model := npc.NPCModel{}

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

// найти  в комнате NPC
func (p *PostgresNPCRepository) FindByRoom(roomId string) ([]*npc.NPC, error) {
	model := []*npc.NPCModel{} //слайс моделей

	err := p.db.Where("room_id = ?", roomId).Find(&model).Error
	if err != nil {
		return nil, err
	}

	//превращаем модели в сущности все
	var npcs []*npc.NPC

	for _, model := range model {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		npcs = append(npcs, entity)

	}

	return npcs, nil
}

// удалить торговца (пока ненадо но пусть будет)
func (p *PostgresNPCRepository) Delete(id string) error {

	result := p.db.Delete(&npc.NPCModel{}, "id = ?", id)

	return result.Error
}

// инициализация NPC
func (r *PostgresNPCRepository) initNPCs() error {
	//Проверяем есть ли NPC уже в БД
	var count int64
	if err := r.db.Model(&npc.NPCModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil //npc есть - ничего не делаем
	}

	//создаем npc
	junkTrader := npc.NewJunkTrader()
	weaponTrader := npc.NewWeaponTrader()

	//Сохраняем их ВСЕХ
	if err := r.Save(junkTrader); err != nil {
		return err
	}

	if err := r.Save(weaponTrader); err != nil {
		return err
	}

	return nil

}
