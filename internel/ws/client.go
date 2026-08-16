package ws

import (
	"context"
	"encoding/json"
	"log"
	"project_chat/internel/dto"

	"github.com/coder/websocket"
	"github.com/go-playground/validator/v10"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.CloseNow()
	}()

	for {
		_, raw, err := c.conn.Read(context.Background())
		if err != nil {
			break
		}
		var msg dto.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if err := msg.Validate(); err != nil {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				log.Printf("field %s failed on %s", fieldErr.Field(), fieldErr.Tag())
			}
			continue
		}

		out, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		c.hub.Broadcast(out)
	}
}

func (c *Client) WritePump() {
	defer c.conn.CloseNow()

	for msg := range c.send {
		err := c.conn.Write(context.Background(), websocket.MessageText, msg)
		if err != nil {
			return
		}
	}
}
