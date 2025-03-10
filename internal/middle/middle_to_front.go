package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

type FrontEndFunctions struct {
	friend_refresh func(friends []api.Friend)
	user_refresh   func(users []api.User)
	fr_refresh     func(in []api.User, out []api.User)
	hubs_refresh   func(hubs []api.Hub)
}

var frontend_functions FrontEndFunctions

func Pass(refreshfriends func(friends []api.Friend), refreshusers func(users []api.User), hubs_refresh func(hubs []api.Hub)) {
	frontend_functions.friend_refresh = refreshfriends
	frontend_functions.user_refresh = refreshusers
	frontend_functions.hubs_refresh = hubs_refresh
}

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

func HasOutgoingFriendRequest(userID api.ID) bool {
	_, exists := friend_requests[userID]
	return exists
}

func HasIncomingFriendRequest(userID api.ID) bool {
	_, exists := ext_friend_requests[userID]
	return exists
}

func SendFriendRequest(userID api.ID, accept bool) {
	user, ok := users[userID]
	if !ok {
		if friend, ok := ext_friend_requests[userID]; ok {
			user = friend.User
		} else {
			log.Printf("Error: User %v not in user list.", userID)
			return
		}
	}

	cmgr := connection.GetCMGR()
	// [self] is a package variable, see users.go.
	fr := api.FriendRequest{
		Friend:   personal.GetSelf(),
		Accepted: accept,
	}

	cmgr.SendFriendRequest(user, fr)

	// this has to be different depending on whether we are sending a request or response

	if incomingFriend, ok := ext_friend_requests[userID]; ok {
		// user has to be in ext_friend_requests
		if accept {
			// add to friends
			removeUser(userID)
			addNewFriend(userID, incomingFriend)
		}
		delete(ext_friend_requests, userID)
	} else {
		// sending out a new request
		friend_requests[userID] = user
	}

	frontend_functions.user_refresh(getUserList())
}
