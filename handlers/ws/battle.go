package ws

import (
	"akari/models/battle"
)

func OnSearch(c *Client, data interface{}) {
	d := battle.LaunchGame()
	c.Send("launchGame", d)
}
