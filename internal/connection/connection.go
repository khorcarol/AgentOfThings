package connection

import (
	"sync"
)

var (
	cmgr *ConnectionManager
	lock sync.Mutex
)

func GetCMGR() (*ConnectionManager, error) {
	if cmgr == nil {
		lock.Lock()
		defer lock.Unlock()
		if cmgr == nil {
			_cmgr, err := InitConnectionManager()
			if err != nil {
				return nil, err
			}
			cmgr = _cmgr
		}
	}
	return cmgr, nil
}

// var (
// 	// B->M, sends a new discovered user
// 	IncomingUsers chan api.User = make(chan api.User)
//
// 	// B->M, sends the response to the fried request (potentially with data)
// 	IncomingFriendResponse chan api.FriendResponse = make(chan api.FriendResponse)
//
// 	// B->M, sends an external friend request to respond to
// 	IncomingFriendRequest chan api.Friend = make(chan api.Friend)
// )
