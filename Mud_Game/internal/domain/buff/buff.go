package buff

import "time"

// тип бафа
type Type string

const (
	HealthRegen    Type = "health_regen"     //реген ХП
	MaxHealthBoost Type = "max_health_boost" //+макс хп
)

// представляет временный эффект
type Buff struct {
	ID            string        //уникальный ID бафа
	Type          Type          //тип бафа
	Value         int           //значение (например 1 реген/сек)
	Interval      time.Duration //интервал срабатывания(0 = одноразовый/мгновенный)
	Duration      time.Duration //длительность действия
	RemainingTime time.Duration //сколько осталось времени
	Description   string        //описание бафа
	LastTick      time.Time     //время последнего срабатывания
}
