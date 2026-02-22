package tcp

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/pkg/logger"
	"fmt"
	"net"
	"strconv"
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
	conn.Write([]byte("Добро пожаловать в MUD игру! Как тебя зовут?\n> "))
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
		s.logger.Error("Не удалось сохранить игрока " + err.Error())                //это вылезет мне в терминале как админу
		conn.Write([]byte("Ошибка во время создания персонажа. Попробуй позже.\n")) //это отправится игроку
		return
	}

	s.logger.Info("Новый игрок :" + name + "(ID:" + id + ")")
	conn.Write([]byte("Привет " + name + "! Добро пожаловать в игру !\n> "))
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
		////////////////////////////////////// команды ://///////////////////////////////////
		if comand == "" {
			conn.Write([]byte("Введите команду\n> "))
			continue
		}
		//ПОКИНУТЬ ПРИЛОЖЕНИЕ
		if comand == "quit" {
			s.logger.Info("Игрок выходит и удаляется из репозитория")
			conn.Write([]byte("До свидания!\n"))
			s.playerRepo.Delete(newPlayer.ID)
			break
		}
		if comand == "look" { //ОСМОТРЕТЬСЯ
			currentRoomID := newPlayer.CurrentRoom          // 1. Получить текущую комнату игрока(id комнаты)
			room, err := s.roomRepo.FindByID(currentRoomID) //Обращаемся к репозиторию комнат (s.roomRepo) и просим найти комнату по этому ID.
			if err != nil {
				conn.Write([]byte("Комната  не найдена\n> "))
				continue //переходим к следующей итерации цикла, ждём новую команду
			}
			responce := room.Look(newPlayer.ID)   //Получаем описание комнаты
			conn.Write([]byte(responce + "\n> ")) //Добавляем приглашение \n> для следующей команды

			continue
		}
		//ДВИЖЕНИЯ "MOVE"
		if strings.HasPrefix(comand, "move") { //проверяет, начинается ли команда с "move "
			direction := strings.TrimPrefix(comand, "move ")        //убирает "move " и возвращает направление
			room, err := s.roomRepo.FindByID(newPlayer.CurrentRoom) //Обращаемся к репозиторию комнат (s.roomRepo) и просим найти комнату по этому ID.
			if err != nil {
				conn.Write([]byte("Комната  не найдена\n> "))
				continue
			}
			exits := room.GetExits()           //получить карту выходов
			nextRoomID, ok := exits[direction] //Проверить, есть ли такое направление(direction)
			if !ok {
				conn.Write([]byte("Туда нельзя идти\n> "))
				continue
			}
			newPlayer.CurrentRoom = nextRoomID //Обновить позицию игрока и сохранить
			s.playerRepo.Save(newPlayer)

			nextRoom, _ := s.roomRepo.FindByID(nextRoomID) //Показываем новую комнату
			conn.Write([]byte(nextRoom.Look(newPlayer.ID) + "\n> "))
			continue
		}
		//ИНВЕНТАРЬ
		if comand == "inventory" {
			invent := newPlayer.Inventory //смотрим в инвентарь
			if len(invent) == 0 {
				conn.Write([]byte("Инвентарь пуст\n> "))
				continue
			}
			//Считаем сколько каких предметов
			itemCounts := make(map[string]int)
			for _, item := range newPlayer.Inventory {
				itemCounts[item]++
			}
			//Красивый вывод
			var builder strings.Builder
			builder.WriteString("Твой инвентарь:\n")

			for item, count := range itemCounts {
				builder.WriteString(" • ")
				builder.WriteString(item)
				if count > 1 {
					builder.WriteString(" x")
					builder.WriteString(strconv.Itoa(count))
				}
				builder.WriteString("\n")
			}
			builder.WriteString("> ")
			conn.Write([]byte(builder.String()))
			continue

		}
		//ВЗЯТЬ
		if strings.HasPrefix(comand, "take ") {
			itemName := strings.TrimPrefix(comand, "take ") //узнали название предмета

			r, err := s.roomRepo.FindByID(newPlayer.CurrentRoom) // узнали в какой сейчас комнате
			if err != nil {
				conn.Write([]byte("Ошибка загрузки комнаты\n> "))
				continue
			}
			// Вся логика поиска и удаления — ВНУТРИ комнаты!
			takeItem, err := r.TakeItem(itemName)
			if err != nil {
				conn.Write([]byte("Здесь нет такого предмета\n> "))
				continue
			}
			// Добавляем в инвентарь игрока и сохраняем изменения
			newPlayer.Inventory = append(newPlayer.Inventory, takeItem)
			s.playerRepo.Save(newPlayer)

			//Сохраняем изменения в комнате
			s.roomRepo.Save(r)

			conn.Write([]byte("Ты взял:" + takeItem + "\n> "))
			continue

		}

		// Формируем неизвестный ответ
		responce := "Вы ввели неизвестную команду\n> "
		// Отправляем ответ
		conn.Write([]byte(responce))

	}
}
