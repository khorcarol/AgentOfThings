package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"

	priorityQueue "github.com/khorcarol/AgentOfThings/lib/priorityQueue"
)

var users = make(map[string]api.User)
var friends = make(map[string]api.User)

// Modify scoreUser to score a pair of users instead of a single user
func scoreUserPair(user1 api.User, user2 api.User) int {
	score := 0
	// Count common interests between both users
	interests1 := make(map[api.Interest]bool)

	for _, interest := range user1.CommonInterests {
		interests1[interest] = true
	}

	for _, interest := range user2.CommonInterests {
		if interests1[interest] {
			score++
		}
	}

	// Penalise if either user has been seen
	if user1.Seen || user2.Seen {
		score -= 1
	}

	return score
}

// Replace rankUsers with matchUsers
func matchUsers() []priorityQueue.Pair[api.User] {
	// Create a priority queue to store user pairs
	totalPairs := (len(users) * (len(users) - 1)) / 2
	pq := priorityQueue.NewPriorityQueue[priorityQueue.Pair[api.User]](totalPairs)

	// Compare each pair of users
	for id1, user1 := range users {
		for id2, user2 := range users {
			// Skip if same user or if either user is already a friend
			if id1 >= id2 {
				continue
			}
			if _, ok := friends[id1]; ok {
				continue
			}
			if _, ok := friends[id2]; ok {
				continue
			}

			score := scoreUserPair(user1, user2)
			pair := priorityQueue.Pair[api.User]{
				First:  user1,
				Second: user2,
				Score:  score,
			}
			pq.Push(pair, score)
		}
	}

	return pq.To_pairs()
}

func SetUserSeen(id string, val bool) {
	u, t := users[id]
	if t {
		u.Seen = val
		users[id] = u
	}
}
