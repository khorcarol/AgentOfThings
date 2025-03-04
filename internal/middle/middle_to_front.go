package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

func CommonInterests(userID api.ID) []api.Interest {
	return common_interests[userID] 
}

// A collection of functions to be used by the front end
func Seen(userID api.ID) {
	setUserSeen(userID, true)
}


func GetFriends() []api.Friend {
	ret := []api.Friend{}

	for _, v := range friends {
		ret = append(ret, v)
	}

	return ret
}

func SendFriendRequest(userID api.ID) {
	user, err := users[userID]
	if !err {
		log.Printf("Error: Friend response not in ext_friend_requests %t", err)
	}

	cmgr, err2 := connection.GetCMGR()
	if err2 != nil {
		log.Fatal(err)
	}
	// [self] is a package variable, see users.go.
	cmgr.SendFriendRequest(userID, personal.GetSelf())

	delete(users, userID)
	ranked_users.Remove(userID)
	friend_requests[userID] = user
}

// Respond to external friend request
func ExtFriendResponse(userID api.ID, accept bool) {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Printf("Error: Friend response not in ext_friend_requests %b", err)
	}

	cmgr.SendFriendResponse(userID, accept)

	if accept {
		friend, err2 := ext_friend_requests[userID]
		if !err2{
			log.Fatal(err2)
		}
		delete(users, userID)
		ranked_users.Remove(userID)
		friends[userID] = friend 
	}
}
