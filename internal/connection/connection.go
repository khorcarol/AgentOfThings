package connection

import (
	"log"
	"sync"
)

var (
	cmgr *ConnectionManager
	lock sync.Mutex
)

func GetCMGR() *ConnectionManager {
	if cmgr == nil {
		lock.Lock()
		defer lock.Unlock()
		if cmgr == nil {
			_cmgr, err := initConnectionManager()
			if err != nil {
				log.Fatal("Failed to create a connction manager.", err)
			}
			cmgr = _cmgr
		}
	}
	return cmgr
}
