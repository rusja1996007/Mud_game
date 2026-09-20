package item

const (
	ColorReset  = "\033[0m"  //стандартный
	ColorGreen  = "\033[32m" //редкий
	ColorRed    = "\033[31m" //очень редкий
	ColorYellow = "\033[33m" //для квестовых предметов
)

// ИСПОЛЬЗУЙ(GetColoredName)возвращает цвет предмета в зависимости от редкости
func GetItemColor(stack *ItemStack) string {
	if stack == nil {
		return ""
	}

	switch stack.Rarity {
	case "quest":
		return ColorYellow
	case "epic":
		return ColorRed
	case "rare":
		return ColorGreen
	default:
		return ColorReset
	}
}

// возвращает цветное имя
func GetColoredName(stack *ItemStack) string {
	if stack == nil {
		return ""
	}
	return GetItemColor(stack) + stack.Name + ColorReset
}
