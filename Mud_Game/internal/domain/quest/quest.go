package quest

// что из себя представляет квест
type Quest struct {
	ID          string
	Name        string
	Description string
	NPCID       string
	ExpReward   int      //опыт за завершение
	CoinReward  int      //монеты за завршение
	ItemReward  []string //предметы за завершение
}

// прогрес игрока
type PlayerQuest struct {
	PlayerID   string
	QuestID    string
	IsActive   bool
	IsComplete bool
}

// карта всех квестов
var Quests = map[string]Quest{
	"ring_quest": {
		ID:          "ring_quest",
		Name:        "Кольцо для деда",
		Description: "Расстроенный дед потерял кольцо жены. Гоблин украл его! Верни кольцо деду.",
		NPCID:       "sad_old_man",
		ExpReward:   200,
		CoinReward:  15,
		ItemReward:  []string{}, //ничего
	},
}
