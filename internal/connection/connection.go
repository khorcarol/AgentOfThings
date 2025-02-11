package connection

import "github.com/khorcarol/AgentOfThings/internal/api"

type Config struct {
	X int
}

var (
	NewUserChannel       chan api.User    = make(chan api.User)
	FriendRequestChannel chan api.Friends = make(chan api.Friends)
)

func RequestFriend(user api.User, data api.Friend) (success bool, receiptData api.Friend) {
	// start a new stream with friend request protocol
	return false, ""
}
