package middle

import "log"

func Start() {
	go func() {
		for {
			discoverUser() // blocks on [IncomingUsers] channel from cmgr.
		}
	}()

	go func() {
		for {
			friendResonse() // blocks on [IncomingFriendRequest] channel from cmgr.
		}
	}()

	log.Println("Middle layer processing started")
}
