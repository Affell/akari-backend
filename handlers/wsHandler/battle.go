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
		c.Send("submitGrid", map[string]interface{}{"valid": false})
	}

	valid := battle.CheckSolution(c, jsonData.Grid)
	c.Send("submitGrid", map[string]interface{}{"valid": valid})

	if valid {
		battle.EndGame(c, jsonData.Grid)
	}
}
