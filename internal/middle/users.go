package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	priorityQueue "github.com/khorcarol/AgentOfThings/lib/priorityQueue"
)

var (
	users               = make(map[api.ID]api.User)
	friend_requests     = make(map[api.ID]api.User)
	ext_friend_requests = make(map[api.ID]api.Friend)
	friends             = make(map[api.ID]api.Friend)
)

// Assigns a score to a user, based on number of matches
func scoreUser(user api.User) int {
	score := len(user.Interests)
	if user.Seen {
		score -= 100
	}
	return score
}

// Returns a list of users in order of their score
func rankUsers() []api.User {
	pq := priorityQueue.NewPriorityQueue[api.User]()

	for id, user := range users {
		_, ok := friends[id]
		if ok {
			continue
		}
		score := scoreUser(user)

		pq.Push(user, score)
	}

	return pq.To_list()
}

// Sets the user to seen
func setUserSeen(id api.ID, val bool) {
	u, t := users[id]
	if t {
		u.Seen = val
		users[id] = u
	}
}

// Adds a new user to users
func discoverUser() {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	user := <-cmgr.IncomingUsers
	users[user.UserID] = user
	// TODO: find common interests
}

// Recieve response from (our) sent friend request
func friendResonse() {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	friend_res := <-cmgr.IncomingFriendRequest
	// The refactor here is that friend requests can't be rejected, you can just hang indefinitely.
	_, ok := friend_requests[friend_res.User.UserID]
	if ok {
		friends[friend_res.User.UserID] = friend_res
		delete(friend_requests, friend_res.User.UserID)
		// TODO: Tell user that friend requests have been accepted
	} else {
		ext_friend_requests[friend_res.User.UserID] = friend_res
	}
}

// Recieve a friend request from another user
func extFriendRequest() {
	// TODO: Finish external friend requests
}
