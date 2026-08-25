package tcp

import (
	"Mud_game/Mud_Game/internal/delivery/tcp/handlers"
	"Mud_game/Mud_Game/internal/domain/item"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"fmt"
	"math/rand"
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

	// ✅ ЗАПУСКАЕМ ГЛОБАЛЬНЫЙ ТАЙМЕР респавна предметов у входа в подземелье
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		for range ticker.C {
			entranceRoom, err := s.roomRepo.FindByID("dungeon_entrance_goblins")
			if err == nil && entranceRoom != nil {
				entranceRoom.(*room.Room).RegenerateItems()
				s.roomRepo.Save(entranceRoom)

			}
		}
	}()

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
			if currentPlayer.Stats.Health > 0 {
				currentPlayer.HandleDisconnect(s.playerRepo, s.roomRepo)
			}
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
				Health:     100 + 5*strength,
				Strength:   3,
				Dexterity:  3,
				Intelect:   2,
				Tracking:   6,
				Level:      1,
				Experience: 0,
			},
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

		//Восстановление путешествия(передвижение)
		if currentPlayer.Stats.IsTraveling {
			if time.Now().After(currentPlayer.Stats.TravelEndTime) {
				// Завершаем путешествие при входе
				currentPlayer.CurrentRoom = currentPlayer.Stats.TravelTargetRoom
				currentPlayer.Stats.IsTraveling = false
				currentPlayer.Stats.TravelEndTime = time.Time{}
				currentPlayer.Stats.TravelTargetRoom = ""
				s.playerRepo.Save(currentPlayer)

				room, _ := s.roomRepo.FindByID(currentPlayer.CurrentRoom) //комната где сейчас  персонаж
				fmt.Fprintf(conn, "%s\n> ", room.Look(currentPlayer.ID))
			} else {
				remaining := time.Until(currentPlayer.Stats.TravelEndTime).Round(time.Second)
				fmt.Fprintf(conn, "Ты в пути. Осталось: %v. Команды недоступны.\n> ", remaining)

			}
		}

		//Если есть яд - восстановить
		if currentPlayer.Stats.IsPoisoned && currentPlayer.Stats.PoisonTicks > 0 {
			fmt.Fprintf(conn, "⚠️ Ты всё ещё отравлен! Яд продолжает действовать.\n")
			go currentPlayer.StartPoisonTicker(conn, s.playerRepo)
		}

		//Показ комнаты:
		room, _ := s.roomRepo.FindByID(currentPlayer.CurrentRoom)
		fmt.Fprintf(conn, "%s\n> ", room.Look(currentPlayer.ID))

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

}

func (s *Server) routeCommand(conn net.Conn, cmd string, p *player.Player) bool {
	//всегда можем выйти
	if cmd == "quit" {

		//если выход производится во время боя
		if p.CurrentRoom == "dungeon_goblin" {
			room, err := s.roomRepo.FindByID(p.CurrentRoom)
			if err != nil {
				fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
				return false
			}
			monster := room.GetMonster()

			//если монстр жив-автоматический побег с получением урона
			if monster != nil && monster.IsAlive {
				monsterDamage := monster.MinDamage + rand.Intn(monster.MaxDamage-monster.MinDamage+1)
				defence := p.GetTotalDefence()
				reduction := float64(defence) / (float64(defence) + 100)
				finalDamage := int(float64(monsterDamage) * (1 - reduction))
				if finalDamage <= 0 {
					finalDamage = 1
				}

				p.Stats.Health -= finalDamage
				monster.Health = monster.MaxHealth

				msg := fmt.Sprintf("💨 Перед выходом ты сбегаешь и монстр нанёс %d урона вслед.\n> ", finalDamage)
				p.SendMessage(conn, msg)

				if p.Stats.Health <= 0 {
					p.SendMessage(conn, "💀Ты погиб...\n")
					monster.Health = monster.MaxHealth
					room.SetPlayerOccupantID("")
					s.roomRepo.Save(room)
					p.StopAllTickers()
					s.playerRepo.Delete(p.ID)
					conn.Close()
					return true

				}
				//телепорт
				p.CurrentRoom = room.GetExitRoomID()
				room.SetPlayerOccupantID("")
				p.Stats.IsInDungeon = false
				p.Stats.EnteredDungeonAt = time.Time{}
				s.roomRepo.Save(room)
				s.playerRepo.Save(p)

				handlers.HandleQuit(conn, cmd, p, s.roomRepo, s.playerRepo)
				return true
			}
		}
	}
	//путешествие
	if p.PendingTravel {
		if cmd == "yes" {
			p.PendingTravel = false

			if p.PendingTravelDirection == "south" {
				p.Stats.Hunger -= 10
				p.Stats.Thirst -= 20
				p.Stats.TravelTargetRoom = "global_town"
				fmt.Fprintf(conn, "Ты отправляешься в город. Путь займёт 5 минут.\n> ")
			} else if strings.HasPrefix(p.PendingTravelDirection, "дом ") {
				if p.Zone == nil {
					return false
				}
				p.Stats.Hunger -= 10
				p.Stats.Thirst -= 20

				p.Stats.TravelTargetRoom = p.Zone.RoadID
				fmt.Fprintf(conn, "Ты отправляешься домой. Путь займёт 5 минут.\n> ")
			} else if p.PendingTravelDirection == "dungeon" {
				p.Stats.Hunger -= 5
				p.Stats.Thirst -= 5
				p.Stats.TravelTargetRoom = "dungeon_entrance_goblins"
				fmt.Fprintf(conn, "Ты отправляешься к подземелью. Путь займёт 2 минуты.\n> ")
			}

			p.Stats.IsTraveling = true
			p.Stats.TravelEndTime = time.Now().Add(5 * time.Second) //////////////временно
			s.playerRepo.Save(p)

			//запуск горутины для автоматического завершения
			go func() {
				time.Sleep(time.Until(p.Stats.TravelEndTime))

				//безопасно отправляем сообщение
				p.SendMessage(conn, "\nТы прибыл!\n")

				//обновляем состояние
				p.CurrentRoom = p.Stats.TravelTargetRoom
				p.Stats.IsTraveling = false
				p.Stats.TravelEndTime = time.Time{}
				p.Stats.TravelTargetRoom = ""
				s.playerRepo.Save(p)

				room, _ := s.roomRepo.FindByID(p.CurrentRoom)
				p.SendMessage(conn, room.Look(p.ID)+"\n> ")

			}()

			return false
		} else if cmd == "no" {
			p.PendingTravel = false
			fmt.Fprintf(conn, "Путь отменен\n> ")
			return false
		} else {
			fmt.Fprintf(conn, "Сначала подтверди путешествие командой 'yes' или отмени его командой 'no'\n> ")
			return false
		}
	}

	// Если в путешествии — блокируем все команды
	if p.Stats.IsTraveling {
		if time.Now().After(p.Stats.TravelEndTime) {
			p.CurrentRoom = p.Stats.TravelTargetRoom
			p.Stats.IsTraveling = false
			p.Stats.TravelEndTime = time.Time{}
			p.Stats.TravelTargetRoom = ""
			s.playerRepo.Save(p)

			room, err := s.roomRepo.FindByID(p.CurrentRoom)
			if err != nil {
				fmt.Fprintf(conn, "Ошибка загрузки комнаты.\n> ")
				return false
			}

			fmt.Fprintf(conn, "Ты прибыл!\n")
			fmt.Fprintf(conn, "%s\n> ", room.Look(p.ID))
			return false

		} else {

			remaining := time.Until(p.Stats.TravelEndTime).Round(time.Second)
			fmt.Fprintf(conn, "Ты в пути. Осталось: %v. Команды недоступны.\n> ", remaining)
			return false
		}
	}
	//если спишь, блокируем все
	if p.Stats.IsSleeping && cmd != "wake" {
		fmt.Fprintf(conn, "Ты спишь, проснись командой 'wake'.\n> ")
		return false
	}

	//если в отеле то ждем
	if p.Stats.IsSleepingHotel {
		fmt.Fprintf(conn, "Ты отдыхаешь в отеле. Подожди окончания отдыха.\n> ")
		return false
	}

	// Если игрок на охоте — блокируем все команды кроме "hunt"
	if p.Stats.IsHunting {
		if cmd == "hunt" {
			fmt.Fprintf(conn, "Ты на охоте, вернешься через %v\n> ",
				time.Until(p.Stats.HuntingEndTime).Round(time.Second))

		} else if cmd == "quit" {
			handlers.HandleQuit(conn, cmd, p, s.roomRepo, s.playerRepo)
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
	case strings.HasPrefix(cmd, "use "):
		handlers.HandleUse(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "look"):
		handlers.HandleLook(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "pay 20"):
		handlers.HandleHotel(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "attack"):
		handlers.HandleAttack(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "search"):
		handlers.HandleSearch(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case strings.HasPrefix(cmd, "escape"):
		handlers.HandleEscape(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	default:
		fmt.Fprintf(conn, "Неизвестная команда\n> ")
		return false
	}

}
