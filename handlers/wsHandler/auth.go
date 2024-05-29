package wsHandler

import (
	"akari/models/user"
	"akari/models/ws"
	"encoding/json"
)

func OnAuth(c *ws.Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	jsonData := struct {
		Token string `json:"token"`
	}{}
	json.Unmarshal(jsonString, &jsonData)

	if jsonData.Token == "" {
		c.Socket.Close()
		return
	}

	u := user.UserToken{}
	var err error
	if u, err = user.GetUserToken(jsonData.Token); err != nil {
		c.Send("close", nil)
		c.Socket.Close()
		return
	}

	if old, ok := ws.WsUsers[u.ID]; ok && old.Socket != c.Socket {
		old.Send("authAttempt", nil)
		c.Send("close", nil)
		c.Socket.Close()
	}

	c.User = u
	ws.WsUsers[u.ID] = c

	c.Send("authenticated", nil)
}
