package player

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB // ← храним соединение с БД
}

func NewPostgresRepository(db *gorm.DB) (*PostgresRepository, error) {
	if err := db.AutoMigrate(&player.PlayerModel{}); err != nil { //волшебная палочка, создающая таблицы под твои структуры(player.PlayerModel)
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}
func (r *PostgresRepository) Save(p *player.Player) error {

	model, err := player.FromEntity(p)
	if err != nil {

		return err
	}

	result := r.db.Save(model)
	if result.Error != nil {

		return result.Error
	}

	return nil
}

func (r *PostgresRepository) FindByID(id string) (*player.Player, error) {
	model := player.PlayerModel{}
	// Ищем в БД
	//&model — куда положить результат (как коробка)
	//"id = ?" — условие поиска (ищем по полю id)
	//id — значение, которое подставится вместо ?
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // игрок не найден
		}
		return nil, err // другая ошибка

	}
	// Конвертируем модель в сущность
	playerEntity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}
	return playerEntity, nil
}
func (r *PostgresRepository) FindByName(name string) (*player.Player, error) {

	model := player.PlayerModel{}

	err := r.db.Unscoped().Where("name = ?", name).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if model.Deleted_at.Valid {
		return nil, nil //// игрок помечен как удалённый → считаем что его нет
	}

	return model.ToEntity()
}

func (r *PostgresRepository) Delete(id string) error {
	if r.db == nil {

		return errors.New("db is nil")
	}

	result := r.db.Exec("DELETE FROM player_models WHERE id = ?", id)

	return result.Error
}
