package websocket

import (
	"Mud_game/Mud_Game/internal/domain/player"
	"Mud_game/Mud_Game/internal/domain/room"
	"Mud_game/Mud_Game/internal/repository/npc_repo"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	//1-создаем адаптер
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

	//3-ищем игрока
	existingPlayer, err := s.playerRepo.FindByName(name)
	if err != nil {
		adapter.Write([]byte("Ошибка при входе\n"))
		return
	}

	//4-отвечаем
	if existingPlayer != nil {
		fmt.Fprintf(adapter, "С возвращением, %s\n", name)
	} else {
		fmt.Fprintf(adapter, "Привет %s! Добро пожаловать!", name)
	}
}
