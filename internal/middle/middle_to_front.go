package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
)

func CommonInterests(api.User) []api.Interest{
	// TODO: Find common interests
	return make([]api.Interest, 0)
}

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

	data := getPersonalData()
	connection.SendFriendRequest(user, data)

	delete(users, userID)
	friend_requests[userID] = user
}

// Respond to external friend request
func ExtFriendResponse(userID api.ID, accept bool){
	// TODO: Respond with personal data
	// TODO: Get personal data
	//resp := api.FriendResponse{userID, accept, }
	//connection.ExtFriendResponseChannel <- resp
}
