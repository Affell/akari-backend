package wsHandler

import (
	"akari/models/battle"
	"akari/models/grid"
	"akari/models/ws"
	"encoding/json"
)

func OnSearch(c *ws.Client, data interface{}) {
	battle.RegisterPlayer(c)
}

func OnCancelSearch(c *ws.Client, data interface{}) {
	battle.CancelMatchmaking(c)
}

func OnGridSubmit(c *ws.Client, data interface{}) {

	jsonString, _ := json.Marshal(data)
	jsonData := struct {
		Grid grid.Grid `json:"grid"`
	}{}
	json.Unmarshal(jsonString, &jsonData)

	if len(jsonData.Grid) == 0 {
		c.Send("gridSubmit", map[string]interface{}{"valid": false})
		return
	}

	valid := battle.CheckSolution(c, jsonData.Grid)
	c.Send("gridSubmit", map[string]interface{}{"valid": valid})

	if valid {
		battle.EndGame(c, jsonData.Grid)
	}
}

func OnScoreboard(c *ws.Client, data interface{}) {
	jsonString, _ := json.Marshal(data)
	jsonData := struct {
		Offset int `json:"offset"`
	}{}
	json.Unmarshal(jsonString, &jsonData)

	if jsonData.Offset < 0 {
		jsonData.Offset = 0
	}

	scoreboard := battle.GetScoreboard(jsonData.Offset)
	c.Send("scoreboard", scoreboard)
}
