package handlers

import (
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/npc"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// показ всех npc(cmd - allnpc)
func HandleAllNPC(conn net.Conn, npcRepo *npc_repo.PostgresNPCRepository, p *player.Player) {

	npcs, err := npcRepo.FindByRoom(p.CurrentRoom)
	if err != nil {

		fmt.Fprintf(conn, "ОШибка поиска npc\n> ")
		return
	}

	if len(npcs) == 0 {
		fmt.Fprintf(conn, "В комнате нет npc\n> ")
		return
	}

	fmt.Fprintf(conn, "🏪 NPC в комнате:\n\n")
	for _, npc := range npcs {
		fmt.Fprintf(conn, "        %s\n", npc.Name)
		fmt.Fprintf(conn, "  %s\n", npc.Description)
		fmt.Fprintf(conn, "----------------------\n")

	}
	fmt.Fprintf(conn, "\nИспользуйте `talk <имя>` чтобы поговорить\n> ")
}

// начать диалог с ..(cmd - talk)
func HandleTalk(conn net.Conn, cmd string, p *player.Player, npcRepo *npc_repo.PostgresNPCRepository) {
	if cmd == "talk" {
		fmt.Fprintf(conn, "С кем разговаривать? Использование: talk <name>\n> ")
		return
	}

	arg := strings.Fields(cmd)

	if len(arg) < 2 {
		fmt.Fprintf(conn, "С кем разговаривать? Использование: talk <name>\n> ")
		return
	}

	if p.IsTalkin == true {
		fmt.Fprintf(conn, "Ты уже разговариваешь! Используй `stop talk` чтобы прекратить беседу\n> ")
		return
	}

	thatNPC := arg[1]

	npcs, err := npcRepo.FindByRoom(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка поиска NPC\n> ")
		return
	}

	var target *npc.NPC
	for _, n := range npcs {
		if strings.EqualFold(n.Name, thatNPC) {
			target = n
			break
		}

	}

	if target == nil {
		fmt.Fprintf(conn, "NPC %s не найден в комнате\n> ", thatNPC)
		return
	}

	if len(target.Inventory) == 0 {
		fmt.Fprintf(conn, "%s не продает товары\n> ", target.Name)
		return
	}

	//если прошло время - обновить товар?
	if time.Since(target.LastRefresh) >= target.RefreshTime {
		switch target.ID {
		case "junk_trader":
			target.Inventory = npc.GenerateJunkItems()
		case "weapon_trader":
			target.Inventory = npc.GenerateWeaponItems()
		default:
		}

		target.LastRefresh = time.Now()
		npcRepo.Save(target)
	}
	//показ товаров и установка флага диалога
	p.IsTalkin = true
	p.TalkToID = target.ID
	p.TalkToName = target.Name
	fmt.Fprintf(conn, "🧙 %s:\n", target.Name)
	fmt.Fprintf(conn, "%s\n", target.Description)
	fmt.Fprintf(conn, "📦 Товары:\n")

	for i, item := range target.Inventory {
		fmt.Fprintf(conn, "%d. %s - %d монет (в наличии: %d)\n", i+1, item.ItemData.Name, item.Price, item.Count)
	}
	fmt.Fprintf(conn, "\nИспользуй 'buy <номер>' чтобы купить\n")
	fmt.Fprintf(conn, "Используй 'stop talk' чтобы закончить разговор\n> ")

}

// купить (cmd buy ...)
func HandleBuy(cmd string, conn net.Conn, p *player.Player, npcRepo *npc_repo.PostgresNPCRepository) {

	if !p.IsTalkin {
		fmt.Fprintf(conn, "Ты ни с кем не разговариваешь\n> ")
		return
	}

	target, err := npcRepo.FindByID(p.TalkToID)
	if err != nil || target == nil {
		fmt.Fprintf(conn, "NPC не найден\n> ")
		p.IsTalkin = false
		return
	}

	target.Mu.Lock()
	defer target.Mu.Unlock()

	args := strings.Fields(cmd)

	if len(args) < 2 {
		fmt.Fprint(conn, "Что купить? Используй `buy <номер предмета`\n> ")
		return
	}

	itemNum, err := strconv.Atoi(args[1])
	if err != nil || itemNum < 1 {
		fmt.Fprintf(conn, "Неверный номер\n> ")
		return
	}

	if itemNum > len(target.Inventory) {
		fmt.Fprintf(conn, "Товара с номером %d нет\n> ", itemNum)
		return
	}

	thatItem := target.Inventory[itemNum-1]

	if thatItem.Count == 0 {
		fmt.Fprintf(conn, "Товар закончился \n> ")
		return
	}

	if !player.HasItem(p.Inventory, "coin", thatItem.Price) {
		fmt.Fprintf(conn, "Недостаточно монет! Нужно : %d\n> ", thatItem.Price)
		return
	}

	if !p.CanAddItem() {
		fmt.Fprintf(conn, "Нет места в инвентаре\n> ")
		return
	}

	//создаем  сам предмет
	newItem := item.GetItem(thatItem.ItemData.Name, 1)
	if newItem == nil {
		fmt.Fprintf(conn, "Ошибка создания предмета\n> ")
		return
	}

	player.RemoveItem(&p.Inventory, "coin", thatItem.Price)

	if !p.AddItemToInventory(newItem) {
		fmt.Fprintf(conn, "Ошибка добавления предмета\n> ")
		return
	}
	thatItem.Count--
	npcRepo.Save(target)

	fmt.Fprintf(conn, "Ты купил %s за %d монет\n> ", thatItem.ItemData.Name, thatItem.Price)

}

// закончить диалог cmd - stop talk
func HandleStopTalk(conn net.Conn, p *player.Player) {
	if !p.IsTalkin {
		fmt.Fprintf(conn, "Ты ни с кем не разговариваешь\n> ")
		return
	}
	name := p.TalkToName
	p.IsTalkin = false
	p.TalkToID = ""
	p.TalkToName = ""
	fmt.Fprintf(conn, "Ты закончил разговор с %s\n> ", name)
}
