package connection

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)

var (
	// B->M, sends a new discovered user
	NewUserChannel           chan api.User   = make(chan api.User)

	// M->B, sends a friend request to be made
	FriendRequestChannel     chan api.User = make(chan api.User)

	// B->M, sends the response to the fried request (potentially with data)
	FriendResponseChannel    chan api.FriendResponse = make(chan api.FriendResponse)

	// B->M, sends an external friend request to respond to
	ExtFriendRequestChannel  chan api.ID = make(chan api.ID)

	// M->B, sends the response to the external friend request
	ExtFriendResponseChannel chan api.FriendResponse = make(chan api.FriendResponse)
)

func RequestFriend(user api.User, data api.Friend) (success bool, receiptData api.Friend) {
	// start a new stream with friend request protocol
	return false, api.Friend{}
}
