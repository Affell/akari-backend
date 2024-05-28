package battle

import (
	"akari/models/grid"
)

const (
	DEFAULT_SIZE       = 10
	DEFAULT_DIFFICULTY = 1
)

func LaunchGame() map[string]interface{} {

	g := grid.GenerateGrid(DEFAULT_SIZE, DEFAULT_DIFFICULTY)
	return map[string]interface{}{
		"grid":       g,
		"size":       DEFAULT_SIZE,
		"difficulty": DEFAULT_DIFFICULTY,
	}

}
