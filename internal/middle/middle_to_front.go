package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

type FrontEndFunctions struct {
	friend_refresh func(friends []api.Friend)
	user_refresh func(users []api.User)
}

var frontend_functions FrontEndFunctions

func Pass(refreshfriends  func(friends []api.Friend), refreshusers func(users []api.User) ){
	frontend_functions.friend_refresh = refreshfriends
	frontend_functions.user_refresh = refreshusers
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

func SendFriendRequest(userID api.ID) {
	user, err := users[userID]
	if !err {
		log.Printf("Error: User not in user list%t", err)
	}

	cmgr, err2 := connection.GetCMGR()
	if err2 != nil {
		log.Fatal(err)
	}
	// [self] is a package variable, see users.go.
	cmgr.SendFriendRequest(user, personal.GetSelf())

	delete(users, userID)
	ranked_users.Remove(userID)
	friend_requests[userID] = user
}

