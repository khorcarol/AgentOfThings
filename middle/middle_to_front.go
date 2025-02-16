package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/middle"
)


// A collection of functions to be used by the front end
func seen(userID api.ID) {
	middle.SetUserSeen(userID, true)
}

func sendFriendRequest(userID api.ID){

}

