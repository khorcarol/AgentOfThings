package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)

// A collection of functions to be used by the front end
func Seen(userID api.ID) {
	setUserSeen(userID, true)
}

func SendFriendRequest(userID api.ID) {
	user, ok := users[userID]
	if (!ok){
		// TODO: Make Error
		return
	}

	delete(users, userID)
	friend_requests[userID] = user

	// TODO: send user along backend channel
}
