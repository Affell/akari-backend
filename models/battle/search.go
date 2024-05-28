package battle

import (
	"akari/handlers/ws"
)

var queue []ws.Client

func RegisterUser(c ws.Client) {

	// if len(queue) > 0 {
	// 	launchGame(c, queue[len(queue)-1])
	// }
	launchGame(c)

}
