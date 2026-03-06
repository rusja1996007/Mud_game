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
	model, err := player.FromEntity(p) // Конвертируем игрока в модель
	if err != nil {
		return err // возвращаем ошибку конвертации
	}
	if err := r.db.Save(model).Error; err != nil { // 2. Сохраняем в БД и проверяем ошибку
		return err // ошибка сохранения
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
	if err := r.db.First(&model, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //gorm.ErrRecordNotFound - "Запись не найдена" 🔍
			return nil, nil
		}
		return nil, err
	}
	playerEntity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}
	return playerEntity, nil
}
func (r *PostgresRepository) Delete(id string) error {
	//     Просто говорим "удали где id = ?"
	// Если запись есть - удалит
	// Если нет - ничего не сделает (и не ошибка)
	return r.db.Delete(&player.PlayerModel{}, "id = ?", id).Error
}
