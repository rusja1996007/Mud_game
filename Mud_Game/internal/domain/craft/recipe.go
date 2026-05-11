package craft

// представляет рецепт крафта
type Recipe struct {
	ID          string         //уник. ID предмета обязательно если две строки через _
	Result      string         //результат(название предмета)
	Count       int            // кол-во результата
	Ingredients map[string]int //что нужно:название-кол-во
	ExpReward   int            //сколько дает опыта

}

// список всех рецептов
var Recipes = []Recipe{
	{
		ID:     "vegetable_set",
		Result: "vegetable set",
		Count:  1,
		Ingredients: map[string]int{
			"burdock":   2,
			"tomato":    2,
			"potato":    3,
			"empty bag": 1,
		},
		ExpReward: 5,
	},
}
