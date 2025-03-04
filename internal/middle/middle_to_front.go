package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

func CommonInterests(api.User) []api.Interest {
	// TODO: Find common interests
	return make([]api.Interest, 0)
}

// A collection of functions to be used by the front end
func Seen(userID api.ID) {
	setUserSeen(userID, true)
}

func SendFriendRequest(userID api.ID) {
	user, ok := users[userID]
	if !ok {
		// TODO: Make Error
		return
	}

	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	// [self] is a package variable, see users.go.
	cmgr.SendFriendRequest(user, personal.GetSelf())

	delete(users, userID)
	friend_requests[userID] = user
}

// Respond to external friend request
func ExtFriendResponse(userID api.ID, accept bool) {
	// TODO: Respond with personal data
	// TODO: Get personal data
	// resp := api.FriendResponse{userID, accept, }
	// connection.ExtFriendResponseChannel <- resp
}
