package ws

import (
	"context"
	"encoding/json"
	"log"
	"project_chat/internel/dto"
	"project_chat/internel/storage"

	"github.com/coder/websocket"
	"github.com/go-playground/validator/v10"
)

type Client struct {
	hub          *Hub
	conn         *websocket.Conn
	send         chan []byte
	username     string
	user_id      int
	room_id      int
	message_Repo *storage.MessageRepository
}

func NewClient(hub *Hub, conn *websocket.Conn, username string, userid int, roomid int, messageRepo *storage.MessageRepository) *Client {
	return &Client{
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		username:     username,
		user_id:      userid,
		room_id:      roomid,
		message_Repo: messageRepo,
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

		msg.Author = c.username

		if err := msg.Validate(); err != nil {
			if valerrors, ok := err.(validator.ValidationErrors); ok {
				for _, fielderror := range valerrors {

					log.Printf("field %s failed on %s", fielderror.Field(), fielderror.Tag())
				}

			}

			errjson := dto.ErrorMessage{
				Type:    "error",
				Payload: err.Error(),
			}
			out, errclient := json.Marshal(errjson)
			if errclient != nil {
				continue
			}
			c.send <- out

			continue
		}

		if err := c.message_Repo.Save(context.Background(), c.user_id, c.room_id, msg.Payload); err != nil {
			log.Println("failed to save message:", err)

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

func (c *Client) Send(msg []byte) {
	c.send <- msg
}
