package tcp

import (
	"Mud_game/Mud_Game/internal/delivery/common"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	port       string
	logger     logger.Logger
	listener   net.Listener                    //"слушатель"- обьект который принимает пподключение
	playerRepo player.Repository               //Чтобы сервер имел доступ к методам сохранения и поиска игроков
	roomRepo   room.Repository                 //все комнаты
	npcRepo    *npc_repo.PostgresNPCRepository //npc
}

// конструктор
func NewServer(port string, log logger.Logger, repo player.Repository, roomRepo room.Repository, npcRepo *npc_repo.PostgresNPCRepository) *Server {
	return &Server{
		port:       port,
		logger:     log,
		playerRepo: repo,
		roomRepo:   roomRepo,
		npcRepo:    npcRepo,
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

	// Найти или создать игрока (общая логика с WebSocket)
	loadedPlayer, err, welcomeMsg := common.InitPlayer(name, s.playerRepo, s.roomRepo)
	if err != nil {
		s.logger.Error("Ошибка InitPlayer: " + err.Error())
		fmt.Fprintf(conn, "Ошибка при входе в игру\n> ")
		return
	}

	currentPlayer = loadedPlayer

	fmt.Fprintf(conn, "%s> ", welcomeMsg)
	s.logger.Info("Игрок подключился с ником: " + name + " (ID: " + currentPlayer.ID + ")")

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
		s.ShowRoomWithNPC(conn, currentPlayer)

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
			if common.RouteCommand(conn, cmd, currentPlayer, s.playerRepo, s.roomRepo, s.npcRepo, s) { // выходим из цикла если routeCommand вернула true (quit)
				break
			}
		}
	}

}

// показ комнаты и npc
func (s *Server) ShowRoomWithNPC(conn net.Conn, p *player.Player) {

	// Проверка на респавн монстра в данже
	if strings.HasPrefix(p.CurrentRoom, "dungeon_") || p.CurrentRoom == "glubini_room" {
		room, _ := s.roomRepo.FindByID(p.CurrentRoom)
		monster := room.GetMonster()

		if monster != nil && !monster.IsAlive && time.Now().After(monster.RespawnTime) {

			//проверяем не появился ли монстр
			if monster.CheckRespawn() {

				s.roomRepo.Save(room)

				p.CurrentRoom = room.GetExitRoomID()
				s.playerRepo.Save(p)
				// Показываем вход
				room, _ = s.roomRepo.FindByID(p.CurrentRoom)
				fmt.Fprintf(conn, "%s", room.Look(p.ID))
				fmt.Fprintf(conn, "\n> ")
				return
			}

		}
	}
	//показ комнату
	room, err := s.roomRepo.FindByID(p.CurrentRoom)
	if err != nil {
		fmt.Fprintf(conn, "Ошибка загрузки комнаты\n> ")
		return
	}

	fmt.Fprintf(conn, "%s", room.Look(p.ID))

	//показ npc
	npcs, err := s.npcRepo.FindByRoom(p.CurrentRoom)
	if err == nil && len(npcs) > 0 {
		fmt.Fprintf(conn, "\n👥 Ты видишь:")
		for _, npc := range npcs {
			fmt.Fprintf(conn, "\n🧔 %s", npc.Name)
		}

	}
	fmt.Fprintf(conn, "\n> ")

}
