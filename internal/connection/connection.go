package connection

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)

var (
	// B->M, sends a new discovered user
	IncomingUsers    chan api.User   = make(chan api.User)

	// B->M, sends the response to the fried request (potentially with data)
	IncomingFriendResponse chan api.FriendResponse = make(chan api.FriendResponse)

	// B->M, sends an external friend request to respond to
	IncomingFriendRequest chan api.Friend = make(chan api.Friend)
)

func SendFriendRequest(user api.User, data api.Friend) (error) {
	// start a new stream with friend request protocol
	return nil
}

func SendFriendResponse(user api.User, data api.FriendResponse) (error) {
	return nil
}
