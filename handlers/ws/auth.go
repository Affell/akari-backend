package ws

import (
	"akari/models/user"
	"encoding/json"
)

func OnAuth(c *Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	jsonData := struct {
		Token string `json:"token"`
	}{}
	json.Unmarshal(jsonString, &jsonData)

	if jsonData.Token == "" {
		c.socket.Close()
		return
	}

	u := user.UserToken{}
	var err error
	if u, err = user.GetUserToken(jsonData.Token); err != nil {
		c.socket.Close()
		return
	}

	if old, ok := WsUsers[u.ID]; ok && old.socket != c.socket {
		old.Send("close", nil)
		old.socket.Close()
	}

	c.User = u
	WsUsers[u.ID] = c

	c.Send("authenticated", nil)
}
