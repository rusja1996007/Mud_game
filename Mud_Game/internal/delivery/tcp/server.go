package tcp

import (
	"Mud_game/Mud_Game/internal/delivery/tcp/handlers"
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	port       string
	logger     logger.Logger
	listener   net.Listener      //"слушатель"- обьект который принимает пподключение
	playerRepo player.Repository //Чтобы сервер имел доступ к методам сохранения и поиска игроков
	roomRepo   room.Repository   //все комнаты
}

// конструктор
func NewServer(port string, log logger.Logger, repo player.Repository, roomRepo room.Repository) *Server {
	return &Server{
		port:       port,
		logger:     log,
		playerRepo: repo,
		roomRepo:   roomRepo,
		//listenet - nill, создастся позже
	}

}
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port) //Создаем "слушателя" - net.Listen("tcp", ":4000") - говорит ОС: "слушай порт 4000, отдавай мне все подключения"
	if err != nil {
		return err //Если порт занят или нет прав
	}
	// . Сохраняем слушателя в структуру
	s.listener = listener
	s.logger.Info("Запуск сервера TCP по порту :" + s.port)
	for {
		//Ждем подключения (блокируется до появления игрока)
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("ОШибка подключения :" + err.Error())
			continue
		}
		//  Запускаем обработчик в отдельной горутине
		// Каждый игрок работает параллельно!
		go s.handleConnection(conn)
	}
}
func (s *Server) handleConnection(conn net.Conn) { //Метод handleConnection - общение с игроком
	// Гарантированно закрываем соединение при выходе из функции
	defer conn.Close()

	var currentPlayer *player.Player

	defer func() {
		if currentPlayer != nil {
			currentPlayer.StopAllTickers()
		}
	}()
	fmt.Printf("🔌 Новое подключение\n")

	// Отправляем приветствие
	// conn.Write принимает []byte, преобразуем строку в байты и для переноса строки -\n
	fmt.Fprintf(conn, "Добро пожаловать в MUD игру! Как тебя зовут?\n> ")
	// Создаем буфер для чтения команд
	// 1024 байт достаточно для любой команды
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer) // ждем пока напечатает имя
	if err != nil {
		s.logger.Info("Игрок отключился во время ввода имени")
		return
	}
	name := string(buffer[:n]) //преобразованое имя
	name = name[:len(name)-2]  ////преобразованое имя без \r \n в конце

	//Ищем игрока в БД
	existingPlayer, err := s.playerRepo.FindByName(name)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка при входе в игру\n> ")
		return
	}

	if existingPlayer != nil {
		//Игрок найден - загружаем
		currentPlayer = existingPlayer
		fmt.Fprintf(conn, "С возвращением, %s!\n> ", name)
	} else {
		// Новый игрок - создаём
		//Геренириуем ID
		id := fmt.Sprintf("player_%d", time.Now().UnixNano())

		// 1. Создаём личную зону для игрока
		zone := player.CreatePlayerZone(id, name)

		// 2. Сохраняем все комнаты зоны в БД
		for _, room := range zone.Rooms {
			if err := s.roomRepo.Save(room); err != nil {
				s.logger.Error("Не удалось сохранить комнату игрока: " + err.Error())
				fmt.Fprintf(conn, "Ошибка создания комнаты. Попробуй позже\n")
				return
			}
		}

		//Связываем дорогу игрока с городом
		townInterface, err := s.roomRepo.FindByID("global_town")
		if err != nil {
			s.logger.Error("Не удалось найти город")
			fmt.Fprintf(conn, "Ошибка загрузки города, попробуй позже\n")
			return
		}

		//3. Привести интерфейс к конкретному типу *room.Room
		townRoom, ok := townInterface.(*room.Room)
		if !ok {
			s.logger.Error("Не удалось привести интерфейс к конкретному типу\n")
			return
		}

		//4. Создать название выхода
		nameExit := fmt.Sprintf("дом %s", name)

		// 5. Добавить в Exits (для перемещения)
		townRoom.Exits[nameExit] = zone.RoadID

		// 5. ДОБАВИТЬ В TownExits (самое важное!)
		townRoom.TownExits = append(townRoom.TownExits, room.TownExit{
			Name:    nameExit,
			RoomID:  zone.RoadID,
			OwnerID: id,
		})
		//Теперь город "знает", что выход "nameExit" принадлежит игроку с ID "id".
		if err := s.roomRepo.Save(townRoom); err != nil {
			s.logger.Error("Не удалось обновить город: " + err.Error())
		} else {

		}

		//Создаем игрока
		//изза влияния силы на здоровья вывел отдельно в переменную
		strength := 3
		currentPlayer = &player.Player{
			ID:          id,
			Name:        name,
			CurrentRoom: zone.HomeRoomID,     //стартовая комната
			Inventory:   []*item.ItemStack{}, // пустой инвентарь
			Equipment:   &player.Equipment{},
			Stats: &player.Stats{
				MaxSlots:   8,
				Hunger:     100,
				Thirst:     100,
				Health:     50 + 5*strength,
				Strength:   3,
				Dexterity:  3,
				Intelect:   2,
				Tracking:   2,
				Level:      1,
				Experience: 0},
			Zone: zone,
		}

		//Сохраняем в репозиторий
		err = s.playerRepo.Save(currentPlayer)
		if err != nil {

			s.logger.Error("Не удалось сохранить игрока " + err.Error())               //это вылезет мне в терминале как админу
			fmt.Fprintf(conn, "Ошибка во время создания персонажа. Попробуй позже.\n") //это отправится игроку
			return
		}
		fmt.Printf("✅ Игрок УСПЕШНО создан: %s\n", name)

		s.logger.Info("Новый игрок :" + name + "(ID:" + id + ")")
		fmt.Fprintf(conn, "Привет %s! Добро пожаловать в игру!\n> ", name)
	}

	//Восстановление охоты
	if currentPlayer.Stats.IsHunting {
		if time.Now().After(currentPlayer.Stats.HuntingEndTime) {
			currentPlayer.EndHunt(conn, s.playerRepo, s.roomRepo)
			//после завершения охоты запускаем тикеры и показываем комнату
			go currentPlayer.StartHungerTicker(conn, s.playerRepo)
			go currentPlayer.StartThirstTicker(conn, s.playerRepo)
			room, _ := s.roomRepo.FindByID(currentPlayer.CurrentRoom)
			fmt.Fprintf(conn, "%s\n> ", room.Look(currentPlayer.ID))
		} else {
			fmt.Fprintf(conn, "Ты на охоте! Вернешься через %v\n> ",
				time.Until(currentPlayer.Stats.HuntingEndTime).Round(time.Second))

			go func() {
				time.Sleep(time.Until(currentPlayer.Stats.HuntingEndTime))
				currentPlayer.EndHunt(conn, s.playerRepo, s.roomRepo)
			}()
		}
	} else {

		//запускаем тикер отнимания еды и воды
		go currentPlayer.StartHungerTicker(conn, s.playerRepo)
		go currentPlayer.StartThirstTicker(conn, s.playerRepo)
		go currentPlayer.StartBuffTicker(conn, s.playerRepo)
		room, _ := s.roomRepo.FindByID(currentPlayer.CurrentRoom) //комната где сейчас  персонаж
		fmt.Fprintf(conn, "%s\n> ", room.Look(currentPlayer.ID))
	}
	// Цикл обработки команд одного игрока
	for {
		// Читаем команду от игрока
		// n - сколько байт реально прочитали
		n, err := conn.Read(buffer)
		if err != nil {
			s.logger.Info("Игрок отключился")
			break // Выходим из цикла, сработает defer conn.Close()
		}
		// Преобразуем байты в строку
		// buffer[:n] - берем только прочитанные байты (остальной буфер пустой)
		cmd := string(buffer[:n])
		// Обрезаем символы \r\n (нажатие Enter)
		// Например "help\r\n" станет "help"
		cmd = cmd[:len(cmd)-2]
		////////////////////////////////////// команды ://///////////////////////////////////
		if s.routeCommand(conn, cmd, currentPlayer) { // выходим из цикла если routeCommand вернула true (quit)
			break
		}
	}
}

func (s *Server) routeCommand(conn net.Conn, cmd string, p *player.Player) bool {

	if cmd == "quit" {
		handlers.HandleQuit(conn, cmd, p, s.roomRepo, s.playerRepo)
		return true
	}

	//если спишь, блокируем все
	if p.Stats.IsSleeping && cmd != "wake" {
		fmt.Fprintf(conn, "Ты спишь, проснись командой 'wake'.\n> ")
		return false
	}

	// ✅ Если игрок на охоте — блокируем все команды кроме "hunt"
	if p.Stats.IsHunting {
		if cmd == "hunt" {
			fmt.Fprintf(conn, "Ты на охоте, вернешься через %v\n> ",
				time.Until(p.Stats.HuntingEndTime).Round(time.Second))

		} else {
			fmt.Fprintf(conn, "Ты на охоте! Нельзя использовать команды кроме hunt и quit\n> ")
		}
		return false
	}

	//обработка выбора характеристик
	if p.PendingStatChoiсe {
		if time.Now().After(p.PendingStatChoiсeExpiry) {
			p.PendingStatChoiсe = false
			fmt.Fprintf(conn, "Время ожидания истекло, введите повторно 'statpoints'\n> ")
			return false
		}

		if cmd == "1" || cmd == "2" || cmd == "3" || cmd == "4" {
			p.ProcessStatChoice(cmd, conn)
			return false
		} else {
			fmt.Fprintf(conn, "Некорректный ввод. Введите 1,2,3 или 4\n> ")
			return false
		}
	}

	if p.PendingHunt && cmd != "yes" {
		p.PendingHunt = false
		fmt.Fprintf(conn, "Подтверждение охоты отменено.\n> ")
		return false
	}

	if cmd == "" {
		handlers.HandleEmpty(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	}
	switch {
	///////////////////////////////////////////////////////////////////////////////////
	case cmd == "damage":
		p.Stats.Health -= 40
		if p.Stats.Health < 1 {
			p.Stats.Health = 1
		}
		fmt.Fprintf(conn, "Здоровье уменьшено до %d\n> ", p.Stats.Health)
		return false
		//////////////////////////////////////////////////////////////////////
	case cmd == "sleep":
		handlers.HandleSleep(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "wake":
		handlers.HandleWake(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	case cmd == "yes":
		handlers.HandleYesHunt(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "hunt":
		handlers.HandleHunt(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "quit":

		handlers.HandleQuit(conn, cmd, p, s.roomRepo, s.playerRepo)
		return true // сигнал на выход

	case cmd == "inventory":

		handlers.HandleInventory(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "stats":
		handlers.HandleStats(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	case cmd == "statpoints":
		handlers.HandleStatPoints(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "move "): //после move идет еще чтото. аналогично ниже

		handlers.HandleMove(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "take "):

		handlers.HandleTake(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "drop "):

		handlers.HandleDrop(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "destroy "):

		handlers.HandleDestroy(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "garden"):
		handlers.HandleGarden(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "plant "):
		handlers.HandlePlant(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "harvest "):
		handlers.HandleHarvest(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "wear "):
		handlers.HandleWear(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "remove "):
		handlers.HandleRemove(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "eat "):
		handlers.HandleEat(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "drink "):
		handlers.HandleDrink(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "fill "):
		handlers.HandleFill(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "craft":
		handlers.HandleCraft(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	case strings.HasPrefix(cmd, "craft "):
		handlers.HandleCraft(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "look"):
		handlers.HandleLook(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	default:
		fmt.Fprintf(conn, "Неизвестная команда\n> ")
		return false
	}

}
