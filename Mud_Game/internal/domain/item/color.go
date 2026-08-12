package item

const (
	ColorReset = "\033[0m"  //стандартный
	ColorGreen = "\033[32m" //редкий
	ColorRed   = "\033[31m" //очень редкий
)

// (ИСПОЛЬЗУЙ GetColoredName)возвращает цвет предмета в зависимости от редкости
func GetItemColor(stack *ItemStack) string {
	if stack.FireDamage > 0 ||
		stack.PoisonDamage > 0 ||
		stack.MagicDamage > 0 ||
		stack.FireDefence > 0 ||
		stack.PoisonDefence > 0 ||
		stack.MagicDefence > 0 ||
		stack.Name == "vegetable set" ||
		stack.Name == "inonotus obliquus" {
		return ColorGreen
	}
	return ColorReset

}

// возвращает цветное имя
func GetColoredName(stack *ItemStack) string {
	if stack == nil {
		return ""
	}
	return GetItemColor(stack) + stack.Name + ColorReset
}
