package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
)

import priorityQueue "github.com/khorcarol/AgentOfThings/lib/priorityQueue"

var users = make(map[api.ID]api.User)
var friend_requests = make(map[api.ID]api.User)
var ext_friend_requests = make(map[api.ID]api.User)
var friends = make(map[api.ID]api.Friend)

func getPersonalData() api.Friend {
	// TODO: Finish function
	return api.Friend{User: api.User{}, Photo:"", Name:""}
}

// Assigns a score to a user, based on number of matches
func scoreUser(user api.User) int {
	var score = len(user.Interests)
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
func discoverUser(){
	user := <- connection.IncomingUsers
	users[user.UserID] = user
	// TODO: find common interests
}

// Recieve response from (our) sent friend request
func friendResonse(){
	friend_res := <- connection.IncomingFriendResponse
	if friend_res.Accept{
		friends[friend_res.UserID] = friend_res.Data
	} else {
		// TODO: Inform user that friend request has been rejected
	}
	delete(friend_requests, friend_res.UserID)
}

// Recieve a friend request from another user
func extFriendRequest(){
	// TODO: Finish external friend requests
}
