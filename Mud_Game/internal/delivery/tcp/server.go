package tcp

import (
	"Mud_game/Mud_Game/internal/delivery/tcp/handlers"
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

	//Геренириуем ID
	id := fmt.Sprintf("player_%d", time.Now().UnixNano())

	//Создаем игрока
	newPlayer := &player.Player{
		ID:          id,
		Name:        name,
		CurrentRoom: "home_01", //стартовая комната
	}

	//Сохраняем в репозиторий
	err = s.playerRepo.Save(newPlayer)
	if err != nil {
		s.logger.Error("Не удалось сохранить игрока " + err.Error())               //это вылезет мне в терминале как админу
		fmt.Fprintf(conn, "Ошибка во время создания персонажа. Попробуй позже.\n") //это отправится игроку
		return
	}

	s.logger.Info("Новый игрок :" + name + "(ID:" + id + ")")
	fmt.Fprintf(conn, "Привет %s! Добро пожаловать в игру!\n> ", name)
	room, _ := s.roomRepo.FindByID(newPlayer.CurrentRoom) //комната где сейчас  персонаж
	fmt.Fprintf(conn, "%s\n> ", room.Look(newPlayer.ID))

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
		if s.routeCommand(conn, cmd, newPlayer) { // выходим из цикла если routeCommand вернула true (quit)
			break
		}
	}
}

func (s *Server) routeCommand(conn net.Conn, cmd string, p *player.Player) bool {

	if cmd == "" {
		handlers.HandleEmpty(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false
	}
	switch {
	case cmd == "quit":

		handlers.HandleQuit(conn, cmd, p, s.roomRepo, s.playerRepo)
		return true // сигнал на выход

	case cmd == "look":

		handlers.HandleLook(conn, cmd, p, s.roomRepo, s.playerRepo)
		return false

	case cmd == "inventory":

		handlers.HandleInventory(conn, cmd, p, s.roomRepo, s.playerRepo)
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
	default:
		fmt.Fprintf(conn, "Неизвестная команда\n> ")
		return false
	}

}
