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
}

var frontend_functions FrontEndFunctions

func Pass(refreshfriends func(friends []api.Friend), refreshusers func(users []api.User), fr_refresh func(in []api.User, out []api.User)) {
	frontend_functions.friend_refresh = refreshfriends
	frontend_functions.user_refresh = refreshusers
	frontend_functions.fr_refresh = fr_refresh
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
		IsOld:    false,
		Accepted: accept,
	}

	cmgr.SendFriendRequest(user, fr)

	// this has to be different depending on whether we are sending a request or response

	if _, ok := users[userID]; ok {
		// sending out a new request
		removeUser(userID)
		friend_requests[userID] = user

		frontend_functions.fr_refresh(getFriendRequests())
		frontend_functions.user_refresh(getUserList())
	} else {
		// user has to be in ext_friend_requests
		if accept {
			// add to friends
			friends[userID] = ext_friend_requests[userID]
			frontend_functions.friend_refresh(getFriendList())
		} else {
			// discard! return them to users
			addUser(ext_friend_requests[userID].User)
			frontend_functions.user_refresh(getUserList())
		}
		delete(ext_friend_requests, userID)
		frontend_functions.fr_refresh(getFriendRequests())
	}
}
