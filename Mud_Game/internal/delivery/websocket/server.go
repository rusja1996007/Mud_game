package websocket

import (
	"Mud_game/Mud_Game/internal/delivery/common"
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// WSServer -вебсокет сервер
type WSServer struct {
	port       string
	playerRepo player.Repository
	roomRepo   room.Repository
	npcRepo    *npc_repo.PostgresNPCRepository
}

// Апгрейдер для HTTP → WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// NewServer — создает WebSocket-сервер
func NewServer(port string, playerRepo player.Repository, roomRepo room.Repository, npcRepo *npc_repo.PostgresNPCRepository) *WSServer {
	return &WSServer{
		port:       port,
		playerRepo: playerRepo,
		roomRepo:   roomRepo,
		npcRepo:    npcRepo,
	}
}

func (s *WSServer) Start() error {
	http.HandleFunc("/ws", s.handleConnection)

	fmt.Printf("🚀 WebSocket-сервер запущен на :%s\n", s.port)

	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		return err
	}
	return nil
}

// handleConnection — обрабатывает нового клиента
func (s *WSServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	wSconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}
	defer wSconn.Close()

	//1-создаем адаптер(секретарь)
	adapter := NewConnAdapter(wSconn)

	//2-читаем ник
	buffer := make([]byte, 1024)
	n, err := adapter.Read(buffer)
	if err != nil {
		log.Println("Ошибка чтения ника:", err)
		return
	}
	name := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("📝 Новый клиент: %s\n", name)

	//3-ищем игрока или создаем нового
	currentPlayer, err, welcomeMsg := common.InitPlayer(name, s.playerRepo, s.roomRepo)
	if err != nil {
		log.Println("Ошибка InitPlayer:", err)
		fmt.Fprintf(adapter, "Ошибка при входе\n")
		return
	}

	fmt.Fprintf(adapter, "%s", welcomeMsg)

	//Восстановление охоты
	if currentPlayer.Stats.IsHunting {
		if time.Now().After(currentPlayer.Stats.HuntingEndTime) {
			currentPlayer.EndHunt(adapter, s.playerRepo, s.roomRepo)
			//после завершения охоты запускаем тикеры и показываем комнату
			go currentPlayer.StartHungerTicker(adapter, s.playerRepo)
			go currentPlayer.StartThirstTicker(adapter, s.playerRepo)
			room, _ := s.roomRepo.FindByID(currentPlayer.CurrentRoom)
			fmt.Fprintf(adapter, "%s\n> ", room.Look(currentPlayer.ID))
		} else {
			fmt.Fprintf(adapter, "Ты на охоте! Вернешься через %v\n> ",
				time.Until(currentPlayer.Stats.HuntingEndTime).Round(time.Second))

			go func() {
				time.Sleep(time.Until(currentPlayer.Stats.HuntingEndTime))
				currentPlayer.EndHunt(adapter, s.playerRepo, s.roomRepo)
			}()
		}
	} else {
		//тикеры
		go currentPlayer.StartHungerTicker(adapter, s.playerRepo)
		go currentPlayer.StartThirstTicker(adapter, s.playerRepo)
		go currentPlayer.StartBuffTicker(adapter, s.playerRepo)

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
				fmt.Fprintf(adapter, "%s\n> ", room.Look(currentPlayer.ID))
			} else {
				remaining := time.Until(currentPlayer.Stats.TravelEndTime).Round(time.Second)
				fmt.Fprintf(adapter, "Ты в пути. Осталось: %v. Команды недоступны.\n> ", remaining)

			}
		}

		//Если есть яд - восстановить
		if currentPlayer.Stats.IsPoisoned && currentPlayer.Stats.PoisonTicks > 0 {
			fmt.Fprintf(adapter, "⚠️ Ты всё ещё отравлен! Яд продолжает действовать.\n")
			go currentPlayer.StartPoisonTicker(adapter, s.playerRepo)
		}

		s.ShowRoomWithNPC(adapter, currentPlayer)

		fmt.Printf("📝 Игрок подключился через WS: %s (ID: %s)\n", name, currentPlayer.ID)

		//5 - основной цикл обработки команд:
		for {
			buffer := make([]byte, 1024)
			n, err := adapter.Read(buffer)
			if err != nil {
				log.Printf("Клиент %s отключился: %v\n", name, err)
				break
			}

			cmd := strings.TrimSpace(string(buffer[:n]))
			fmt.Printf("📩 Команда от %s: %s\n", name, cmd)

			//ОБщий роутер команд
			if common.RouteCommand(adapter, cmd, currentPlayer, s.playerRepo, s.roomRepo, s.npcRepo, s) {
				break
			}
		}
	}
}

// показ комнаты и npc
func (s *WSServer) ShowRoomWithNPC(conn net.Conn, p *player.Player) {

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
