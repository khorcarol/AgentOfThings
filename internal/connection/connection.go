package main

type Config struct {
	X int
}

var (
	NewUserChannel       chan string = make(chan string)
	FriendRequestChannel chan string = make(chan string)
)

func RequestFriend(ID string, data string) (success bool, receiptData string) {
	// start a new stream with friend request protocol
	return false, ""
}
