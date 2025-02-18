package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)


// A collection of functions to be used by the front end
func seen(userID api.ID) {
	SetUserSeen(userID, true)
}

func sendFriendRequest(userID api.ID){

}

