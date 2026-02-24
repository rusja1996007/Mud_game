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
		//ОСМОТРЕТЬСЯ
		if comand == "look" {
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
		args, found := strings.CutPrefix(comand, "take ") //CutPrefix проверяет, начинается ли с "take "
		if !found {
			// это не take, просто идём дальше к другим командам
		} else { //Если да → парсим команду
			parts := strings.Fields(args) //разбивает по пробелам
			//args	                  parts
			//"3 bottle"	    ["3", "bottle"]
			//"all bottle"	    ["all", "bottle"]
			//"bottle"	        ["bottle"]
			//"3 big bottle"	["3", "big", "bottle"]

			r, err := s.roomRepo.FindByID(newPlayer.CurrentRoom) // узнали в какой сейчас комнате
			if err != nil {
				conn.Write([]byte("Ошибка загрузки комнаты\n> "))
				continue
			}

			var count int = 1 // сколько предметов брать,по умолчанию берём 1
			var itemName string

			//Парсим → count=3, itemName="Empty bottle"
			// Смотрим, что нам прислали

			if len(parts) == 1 && parts[0] == "all" { //// Это команда "take all" — взять всё из комнаты
				allItems := r.GetItems()         //получаем все предметы из комнаты
				for _, stack := range allItems { //проходим по каждому предмету
					for i := 0; i < stack.Count; i++ {
						takenItem, _ := r.TakeItem(stack.Name)
						newPlayer.Inventory = append(newPlayer.Inventory, takenItem)
					}
				}
				s.playerRepo.Save(newPlayer)
				s.roomRepo.Save(r)
				conn.Write([]byte("Вы взяли все из комнаты\n> "))
				continue

			} else if len(parts) == 1 {
				itemName = parts[0]
			} else {
				//Случай Б: 2 или больше частей
				if parts[0] == "all" { //Если первое слово — "all"
					count = -1                              //специальное значение "все"
					itemName = strings.Join(parts[1:], " ") //itemName = strings.Join(["big", "bottle"], " ") = "big bottle"
				} else {
					num, err := strconv.Atoi(parts[0]) // ← пробуем распарсить число
					if err == nil {                    //значит это число
						count = num
						itemName = strings.Join(parts[1:], " ")
					} else { // это не число и не "all" — значит, название из нескольких слов
						itemName = strings.Join(parts, " ")
					}
				}
			}

			//  Если это не "take all", обрабатываем обычный take
			if itemName != "" {
				items := r.GetItems()
				foundIndex := -1
				for i, stack := range items {
					if stack.Name == itemName {
						foundIndex = i
						break
					}
				}
				if foundIndex == -1 {
					conn.Write([]byte("Здесь нет такого предмета\n> "))
					continue
				}
				// Смотрим, сколько штук этого предмета лежит в комнате.
				available := items[foundIndex].Count
				// Сколько будем брать
				takeCount := count
				if count == -1 {
					takeCount = available //если "all" — берём всё
				}
				if takeCount > available {
					takeCount = available // нельзя взять больше, чем есть
				}
				if takeCount == 0 {
					conn.Write([]byte("Нечего брать\n> "))
					continue
				}

				//берем предметы по одному
				for i := 0; i < takeCount; i++ {
					//// Добавляем в инвентарь игрока и сохраняем изменения+Сохраняем изменения в комнате
					takeItem, _ := r.TakeItem(itemName)
					newPlayer.Inventory = append(newPlayer.Inventory, takeItem)
				}
				s.playerRepo.Save(newPlayer)
				s.roomRepo.Save(r)

				conn.Write([]byte(fmt.Sprintf("Ты взял %d %s\n> ", takeCount, itemName)))
				continue
			}

		}
		//УНИЧТОЖИТЬ
		argss, found := strings.CutPrefix(comand, "destroy ")
		if !found {
			//это не destroy - идем дальше
		} else {
			parts := strings.Fields(argss) //парсим(разбиваем) команду
			var count int = 1
			var itemName string

			if len(parts) == 1 {
				itemName = parts[0]
			} else if len(parts) >= 2 {
				if parts[0] == "all" {
					count = -1
					itemName = strings.Join(parts[1:], " ")
				} else {
					num, err := strconv.Atoi(parts[0])
					if err == nil {
						count = num
						itemName = strings.Join(parts[1:], " ")
					} else {
						itemName = strings.Join(parts, " ")
					}

				}
			}
			//// Сначала считаем, сколько таких предметов есть
			available2 := 0
			for _, item := range newPlayer.Inventory {
				if item == itemName {
					available2++
				}
			}
			if available2 == 0 {
				conn.Write([]byte("У тебя нет такого предмета\n> "))
				continue
			}
			//Определить, сколько уничтожать
			destroyCount := count
			if count == -1 {
				destroyCount = available2
			}
			if destroyCount > available2 {
				destroyCount = available2
			}
			if destroyCount == 0 {
				conn.Write([]byte("Нечего уничтожать\n> "))
				continue
			}
			//    Уничтожить предметы (удалить из инвентаря)
			var newInventory []string
			removed := 0
			for _, item := range newPlayer.Inventory {
				if item == itemName && removed < destroyCount {
					removed++ // пропускаем (уничтожаем) ?
				} else {
					newInventory = append(newInventory, item)
				}
			}
			newPlayer.Inventory = newInventory
			s.playerRepo.Save(newPlayer)
			conn.Write([]byte(fmt.Sprintf("Ты уничтожил %d %s\n> ", destroyCount, itemName)))
			continue
		}

		// Формируем неизвестный ответ
		responce := "Вы ввели неизвестную команду\n> "
		// Отправляем ответ
		conn.Write([]byte(responce))

	}
}
