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

func getUserList() []api.User {
	res := []api.User{}
	for _, e := range ranked_users.To_list() {
		res = append(res, users[e])
	}
	return res
}

func getFriendList() []api.Friend {
	res := []api.Friend{}
	for _, v := range friends {
		res = append(res, v)
	}
	return res
}

func getFriendRequests() ([]api.User, []api.User) {
	in := []api.User{}
	out := []api.User{}

	for _, v := range friend_requests {
		out = append(out, v)
	}

	for _, v := range ext_friend_requests {
		in = append(in, v.User)
	}
	return in, out
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
		if _, ok := users[user.UserID]; !ok {
			users[user.UserID] = user

			updateCommonInterests(user.UserID, user.Interests)
			ranked_users.Push(user.UserID, scoreUser(user))

			frontend_functions.user_refresh(getUserList())
		}
}
// Recieve response from (our) sent friend request
func waitOnFriendRequest() {
	cmgr, err := connection.GetCMGR()
	if err != nil {
		log.Fatal(err)
	}
	friend_res := <-cmgr.IncomingFriendRequest

	// check whether we are receiving a response or request
	// TODO: Handle OLD requests/responses
	_, ok := friend_requests[friend_res.Friend.User.UserID]
	if ok {
		// this is a response to a friend request that we sent out
		// check if it is acceptance or rejection
		if friend_res.Accepted {
			// if accepted, set as friend and remove from requests
			friends[friend_res.Friend.User.UserID] = friend_res.Friend
			frontend_functions.friend_refresh(getFriendList())
		} else {
			// if rejected, set as user and remove from requests
			users[friend_res.Friend.User.UserID] = ext_friend_requests[friend_res.Friend.User.UserID].User
			ranked_users.Push(friend_res.Friend.User.UserID, scoreUser(friend_res.Friend.User))
			frontend_functions.user_refresh(getUserList())
		}
		frontend_functions.fr_refresh(getFriendRequests())
		delete(friend_requests, friend_res.Friend.User.UserID)
	} else {
		// this is a new incoming friend request
		delete(users, friend_res.Friend.User.UserID)
		ranked_users.Remove(friend_res.Friend.User.UserID)

		ext_friend_requests[friend_res.Friend.User.UserID] = friend_res.Friend

		frontend_functions.fr_refresh(getFriendRequests())
		frontend_functions.user_refresh(getUserList())
	}
}
