package middle

import (
	"log"
	"sort"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/personal"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

var (
	users               = make(map[api.ID]api.User)
	friend_requests     = make(map[api.ID]api.User)
	ext_friend_requests = make(map[api.ID]api.Friend)
	friends             = make(map[api.ID]api.Friend)
	common_interests    = make(map[api.ID]([]api.Interest))

	ranked_users = [](struct {
		id    api.ID
		score int
	}){}
)

// Retrieve friends from storage
func init() {
	loadedFriends, err := storage.LoadFriends()
	if err == nil {
		friends = loadedFriends
	}
}

func saveFriends() {
	_ = storage.SaveFriends(friends)
}

func AddFriend(id api.ID, user api.Friend) {
	friends[id] = user
	saveFriends()
}

// Assigns a score to a user, based on number of matches
func scoreUser(user api.User) int {
	score := len(user.Interests)
	if user.Seen {
		score -= 100
	}
	return score
}

// Insert into sorted array of users
func addUser(user api.User) {
	user_score := scoreUser(user)
	i := sort.Search(len(ranked_users), func(idx int) bool { return ranked_users[idx].score >= user_score })
	ranked_users = append(ranked_users, struct {
		id    api.ID
		score int
	}{user.UserID, 0}) // Uses temp of user.UserID as it is easier than making a nil id
	copy(ranked_users[i+1:], ranked_users[i:])
	ranked_users[i] = struct {
		id    api.ID
		score int
	}{user.UserID, user_score}

	users[user.UserID] = user
}

func removeUser(id api.ID) {
	delete(users, id)

	i := sort.Search(len(ranked_users), func(idx int) bool { return ranked_users[idx].id == id })
	if i == len(ranked_users) {
		return
	}
	ranked_users = append(ranked_users[:i], ranked_users[i+1:]...)
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
	for _, e := range ranked_users {
		res = append(res, users[e.id])
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
		addUser(user)
		updateCommonInterests(user.UserID, user.Interests)
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
			addUser(friend_res.Friend.User)
			frontend_functions.user_refresh(getUserList())
		}
		frontend_functions.fr_refresh(getFriendRequests())
		delete(friend_requests, friend_res.Friend.User.UserID)
	} else {
		// this is a new incoming friend request
		removeUser(friend_res.Friend.User.UserID)

		ext_friend_requests[friend_res.Friend.User.UserID] = friend_res.Friend

		frontend_functions.fr_refresh(getFriendRequests())
		frontend_functions.user_refresh(getUserList())
	}
}
