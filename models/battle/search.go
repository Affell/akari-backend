package battle

import "akari/models/ws"

type Queue struct {
	Clients []*ws.Client
}

var queue Queue

func RegisterPlayer(c *ws.Client) {
	if other, ok := queue.Pop(); ok {
		LaunchGame(c, other)
	} else {
		queue.Append(c)
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
