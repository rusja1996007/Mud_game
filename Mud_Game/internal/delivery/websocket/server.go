package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSServer -вебсокет сервер
type WSServer struct {
	port string
}

// Апгрейдер для HTTP → WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// NewServer — создает WebSocket-сервер
func NewServer(port string) *WSServer {
	return &WSServer{
		port: port,
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ WebSocket-клиент подключился!")

	// Пока просто читаем и отправляем обратно

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Клиент отключился")
			return
		}

		fmt.Printf("📩 Получено: %s\n", string(message))
		// Отправляем обратно (эхо для теста)
		if err := conn.WriteMessage(messageType, message); err != nil {
			fmt.Println("Ошибка отправки:", err)
			return
		}
	}
}
