package wserver

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	conn     *websocket.Conn
	username string
	h        *Hub
	send     chan *Message
}

func NewClient(username string,c *websocket.Conn, h *Hub) *Client {
	return &Client{
		conn: c,
		username: username,
		h:    h,
		send: make(chan *Message, 10),
	}
}

func (c *Client) SendMessage() {
	defer func() {
		c.h.Unregister <- c
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			break
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			zap.L().Error("error to marshal message")
		}
		c.conn.WriteMessage(websocket.BinaryMessage, messageJSON)
	}

}

func (c *Client) RecvMessage() {

	defer func() {
		c.h.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			zap.L().Error("read message error", zap.String("error", err.Error()), zap.String("user", c.username))
			break
		}
		msg := &Message{}
		json.Unmarshal(message,msg)

		c.h.Send <- msg

		zap.L().Info("send message", zap.String("user", c.username), zap.String("user_message", string(message)))
	}
}
