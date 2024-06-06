package battle

import (
	"akari/models/grid"
	"akari/models/user"
	"akari/models/ws"
	"time"
)

const (
	DEFAULT_SIZE       = 10
	DEFAULT_DIFFICULTY = 1
)

type Game struct {
	Player1     *ws.Client
	Player2     *ws.Client
	Grid        grid.Grid
	CreatedTime int64
}

var games map[int64]*Game = make(map[int64]*Game)

func LaunchGame(player1, player2 *ws.Client) {

	var game Game

	game.Player1 = player1
	game.Player2 = player2
	games[player1.User.ID] = &game
	games[player2.User.ID] = &game

	g := grid.GenerateGrid(DEFAULT_SIZE, DEFAULT_DIFFICULTY)
	game.Grid = g
	d := map[string]interface{}{
		"grid":       g,
		"size":       DEFAULT_SIZE,
		"difficulty": DEFAULT_DIFFICULTY,
	}

	game.CreatedTime = time.Now().Unix()

	d["opponent"] = player2.User.ToOpponentData()
	player1.Send("launchGame", d)
	d["opponent"] = player1.User.ToOpponentData()
	player2.Send("launchGame", d)
}

func CheckSolution(player *ws.Client, grid grid.Grid) (valid bool) {

	game, ok := games[player.User.ID]
	if !ok {
		return
	}

	return game.Grid.CheckSolution(grid)
}

func EndGame(winner *ws.Client, forfeit bool) {
	game, ok := games[winner.User.ID]
	if !ok {
		return
	}

	other := game.Player1
	if other == winner {
		other = game.Player2
	}

	userOther, err := user.GetUserById(other.User.ID)
	if err != nil {
		return
	}
	userWinner, err := user.GetUserById(winner.User.ID)
	if err != nil {
		return
	}

	newElo1, newElo2 := ComputeResult(userWinner.Score, userOther.Score, 1)
	userWinner.Score = newElo1
	winner.User.Score = newElo1
	userOther.Score = newElo2
	other.User.Score = newElo2

	user.UpdateUser(userWinner, false)
	user.UpdateUser(userOther, false)

	winner.Send("gameResult", map[string]interface{}{
		"result":   "win",
		"forfeit":  forfeit,
		"newElo":   newElo1,
		"eloDelta": newElo1 - winner.User.Score,
	})

	other.Send("gameResult", map[string]interface{}{
		"result":   "defeat",
		"forfeit":  forfeit,
		"newElo":   newElo2,
		"eloDelta": newElo2 - other.User.Score,
	})

	delete(games, winner.User.ID)
	delete(games, other.User.ID)

}

func Forfeit(player *ws.Client) {
	game, ok := games[player.User.ID]
	if !ok {
		return
	}

	other := game.Player1
	if other == player {
		other = game.Player2
	}

	EndGame(other, true)
}
