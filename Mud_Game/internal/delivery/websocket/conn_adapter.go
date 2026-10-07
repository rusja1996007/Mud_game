package websocket

import (
	"bytes"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// ConnAdapter — обертка над *websocket.Conn, реализующая net.Conn
type ConnAdapter struct {
	wsConn *websocket.Conn
	buffer bytes.Buffer
}

// NewConnAdapter — создает новый адаптер
func NewConnAdapter(wsConn *websocket.Conn) *ConnAdapter {
	return &ConnAdapter{
		wsConn: wsConn,
	}
}

// wsAddr — заглушка для net.Addr
type wsAddr struct{}

func (a wsAddr) Network() string { return "websocket" } //Возвращает "websocket" (тип сети)
func (a wsAddr) String() string  { return "ws-client" }

// =================================== Методы net.Conn =====================

// Read — читает данные из WebSocket в буфер
func (c *ConnAdapter) Read(b []byte) (int, error) {
	if c.buffer.Len() == 0 {
		_, message, err := c.wsConn.ReadMessage() //Читаем WS-сообщение
		if err != nil {
			return 0, err
		}
		c.buffer.Write(message)
	}
	return c.buffer.Read(b)

}

// Write — отправляет данные через WebSocket
func (c *ConnAdapter) Write(b []byte) (int, error) {
	err := c.wsConn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Закрывает WebSocket-соединение.( заглушка)
func (c *ConnAdapter) Close() error {
	return c.wsConn.Close()
}

// Возвращают пустой wsAddr для сети -"websocket", для игрока -"ws-client".
func (c *ConnAdapter) LocalAddr() net.Addr {
	return wsAddr{}
}

func (c *ConnAdapter) RemoteAddr() net.Addr {
	return wsAddr{}
}

// net.Conn требует эти методы. Мы не используем дедлайны — просто возвращаем nil (нет ошибки).
func (c *ConnAdapter) SetDeadline(t time.Time) error {
	return nil
}

func (c *ConnAdapter) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *ConnAdapter) SetWriteDeadline(t time.Time) error {
	return nil
}

// для проверки "Убедись, что *ConnAdapter реализует net.Conn"
var _ net.Conn = (*ConnAdapter)(nil)
