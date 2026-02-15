package tcp

import (
	"Mud_game/Mud_Game/internal/pkg/logger"
	"net"
)

type Server struct {
	port     string
	logger   logger.Logger
	listener net.Listener //"слушатель"- обьект который принимает пподключение
}

// конструктор
func NewServer(port string, log logger.Logger) *Server {
	return &Server{
		port:   port,
		logger: log,
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
	conn.Write([]byte("Добро пожаловать в MUD игру!\n"))
	// Создаем буфер для чтения команд
	// 1024 байт достаточно для любой команды
	buffer := make([]byte, 1024)
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
		comand := string(buffer[:n])
		// Обрезаем символы \r\n (нажатие Enter)
		// Например "help\r\n" станет "help"
		comand = comand[:len(comand)-2]
		// Проверяем, не хочет ли игрок выйти
		if comand == "quit" {
			conn.Write([]byte("До свидания!\n"))
			break
		}
		// Формируем ответ
		responce := "Вы написали " + comand + "\n> "
		// Отправляем ответ
		conn.Write([]byte(responce))

	}
}
