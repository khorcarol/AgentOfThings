package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

type FrontEndFunctions struct {
	friendrefresh func(friends []api.Friend)
	userrefresh func(users []api.User)}

var frontendfunctions FrontEndFunctions

func Pass(refreshfriends  func(friends []api.Friend), refreshusers func(users []api.User) ){
	frontendfunctions.friendrefresh = refreshfriends
	frontendfunctions.userrefresh = refreshusers
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

