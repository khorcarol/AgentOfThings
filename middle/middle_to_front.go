package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/middle"
)

// A collection of functions to be used by the front end
func seen(userID api.User) {
	middle.SetUserSeen(userID.UserID, true)
}

func sendFriendRequest(userID api.User) {

}
