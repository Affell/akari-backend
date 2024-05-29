package battle

import (
	"akari/models/ws"

	"github.com/kataras/golog"
)

type Queue struct {
	Clients []*ws.Client
}

var queue Queue

func RegisterPlayer(c *ws.Client) {

	if _, ok := games[c.User.ID]; ok {
		c.Send("search", map[string]interface{}{
			"success": false,
		})
		return
	}

	golog.Debug("Blou", c.User)

	if other, ok := queue.Pop(); ok {
		golog.Debug("Blou1", c.User)
		LaunchGame(c, other)
	} else {
		golog.Debug("Blou2", c.User)
		queue.Append(c)
		c.Send("search", map[string]interface{}{
			"success": true,
		})
	}
}

func CancelMatchmaking(c *ws.Client) {
	queue.Remove(c)
}

func (queue *Queue) Pop() (client *ws.Client, ok bool) {
	if len(queue.Clients) > 0 {
		client = queue.Clients[0]
		queue.Clients = queue.Clients[1:]
		ok = true
	}
	return
}

func (queue *Queue) Append(client *ws.Client) {
	for _, c := range queue.Clients {
		if c == client {
			return
		}
	}
	queue.Clients = append(queue.Clients, client)
}

func (queue *Queue) Remove(client *ws.Client) {
	var new []*ws.Client
	for _, c := range queue.Clients {
		if c != client {
			new = append(new, c)
		}
	}
	queue.Clients = new
}
