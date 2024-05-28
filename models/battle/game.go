package battle

import (
	"akari/handlers/ws"
	"akari/models/grid"
)

const (
	DEFAULT_SIZE       = 10
	DEFAULT_DIFFICULTY = 1
)

func launchGame(player1 ws.Client) {

	g := grid.GenerateGrid(DEFAULT_SIZE, DEFAULT_DIFFICULTY)
	data := map[string]interface{}{
		"grid":       g,
		"size":       DEFAULT_SIZE,
		"difficulty": DEFAULT_DIFFICULTY,
	}

	player1.Send("launchGame", data)
	// player2.Send("launchGame", data)

}
