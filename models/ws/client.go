package ws

import (
	"akari/models/user"

	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type FindHandler func(Event) (HandlerDesc, bool)

type Client struct {
	Socket      *websocket.Conn
	findHandler FindHandler
	User        user.UserToken
}

func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		Socket:      socket,
		findHandler: findHandler,
	}
}

func (c *Client) Send(name string, data interface{}) {
	msg := Message{name, data}
	err := c.Socket.WriteJSON(msg)
	if err != nil {
		golog.Errorf("socket write error: %v\n", err)
	}
}

func (c *Client) Read() {
	var msg Message
	for {
		if err := c.Socket.ReadJSON(&msg); err != nil {
			break
		}
		if handlerDesc, found := c.findHandler(Event(msg.Name)); found && (c.User.ID != 0 || !handlerDesc.AuthRequired) {
			handlerDesc.HandlerFunc(c, msg.Data)
		}
	}

	delete(WsUsers, c.User.ID)

	c.Socket.Close()
}
