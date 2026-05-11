package handlers

import (
	"Mud_game/Mud_Game/internal/domain/craft"
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"fmt"
	"net"
	"strings"
)

func HandleCraft(conn net.Conn, cmd string, p *player.Player, roomRepo room.Repository, playerRepo player.Repository) {
	if cmd == "craft" {
		fmt.Fprintf(conn, "Что скрафтить? Использование: craft <рецепт>\n> ")
		return
	}

	args := strings.Fields(cmd)

	if len(args) < 2 {
		fmt.Fprintf(conn, "Что скрафтить? Использование: craft <рецепт>\n> ")
		return
	}

	recipeId := args[1]      //аргумент после craft
	var recipe *craft.Recipe //рецепт

	//Ищем рецепт
	for _, r := range craft.Recipes {
		if r.ID == recipeId {
			recipe = &r
			break

		}
	}

	if recipe == nil {
		fmt.Fprintf(conn, "Неизвестный рецепт\n> ")
		return
	}

	//проверяем ингредиенты и собираем недостающие

	missing := make([]string, 0)

	for ing, need := range recipe.Ingredients {
		if !player.HasItem(p.Inventory, ing, need) {
			have := 0
			for _, stack := range p.Inventory {
				if stack.Name == ing {
					have = stack.Count
					break
				}
			}
			missing = append(missing, fmt.Sprintf("%s(нужно %d, есть %d)", ing, need, have))

		}
	}

	if len(missing) > 0 {
		fmt.Fprintf(conn, "Не хватает ингредиентов:\n")
		for _, m := range missing {
			fmt.Fprintf(conn, "  • %s\n", m)
		}
		fmt.Fprintf(conn, ">")
		return
	}

	if p.GetFreeSlots() < 1 {
		fmt.Fprintf(conn, "Нет места в инвентаре!\n> ")
		return
	}

	//если есть удаляем предметы
	for ing, need := range recipe.Ingredients {
		player.RemoveItem(&p.Inventory, ing, need)
	}

	//добавляем результат
	resultItem := item.GetItem(recipe.Result, recipe.Count)
	p.AddItemToInventory(resultItem)

	//начисляем опыт
	p.AddExperience(recipe.ExpReward, conn)

	fmt.Fprintf(conn, "Вы скрафтили %s! Получено опыта %d.\n> ", recipe.Result, recipe.ExpReward)

}
