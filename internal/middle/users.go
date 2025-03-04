package middle

import (
	"log"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
	priorityQueue "github.com/khorcarol/AgentOfThings/lib/priorityQueue"
)

var (
	users               = make(map[api.ID]api.User)
	friend_requests     = make(map[api.ID]api.User)
	ext_friend_requests = make(map[api.ID]api.Friend)
	friends             = make(map[api.ID]api.Friend)

	common_interests    = make(map[api.ID]([]api.Interest))
	ranked_users        = priorityQueue.NewPriorityQueue[api.ID]()

)

// Assigns a score to a user, based on number of matches
func scoreUser(user api.User) int {
	score := len(user.Interests)
	if user.Seen {
		score -= 100
	}
	return score
}


// Sets the user to seen
func setUserSeen(id api.ID, val bool) {
	u, t := users[id]
	if t {
		u.Seen = val
		users[id] = u
	}
}


// Updates (or creates) entry for common interests
func updateCommonInterests(userID api.ID, interests []api.Interest) {
	self_interests := personal.GetSelf().User.Interests

	int_map := make(map[api.Interest]bool)
	common := []api.Interest{}

	for _, e1 := range self_interests {
		int_map[e1] = true
	}

	for _, e2 := range interests {
		if int_map[e2] {
			common = append(common, e2)
		}
	}

	common_interests[userID] = common
}


// Adds a new user to users
func discoverUser() {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	user := <-cmgr.IncomingUsers

	// TODO: Check if stored friend

	users[user.UserID] = user

	updateCommonInterests(user.UserID, user.Interests)
	ranked_users.Push(user.UserID, scoreUser(user))
}


// Recieve response from (our) sent friend request
func friendResonse() {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	friend_res := <-cmgr.IncomingFriendResponse
	if friend_res.Accept {
		friends[friend_res.UserID] = friend_res.Data
	} else {
		// TODO: Inform user that friend request has been rejected
	}
	delete(friend_requests, friend_res.UserID)
	ranked_users.Remove(friend_res.UserID)
}


// Recieve a friend request from another user
func extFriendRequest() {
	// TODO: Finish external friend requests
}
