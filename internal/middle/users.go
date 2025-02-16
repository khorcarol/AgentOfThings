package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)

import priorityQueue "github.com/khorcarol/AgentOfThings/lib/priorityQueue"



var users = make(map[api.ID]api.User)
var friends = make(map[api.ID]api.User)

// Assigns a score to a user, based on number of matches
func scoreUser(user api.User) int{
	var score = len(user.CommonInterests);
	if user.Seen{
		score -= 100;
	}
	return score
}

func rankUsers() []api.User{

	pq := priorityQueue.NewPriorityQueue[api.User](len(users))


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

func SetUserSeen(id api.ID, val bool){
	u, t := users[id]
	if t{
		u.Seen = val
		users[id] = u
	}
}
