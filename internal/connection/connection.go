package connection

import "sync"

var (
	cmgr *ConnectionManager
	lock sync.Mutex
)

func GetCMGR() (*ConnectionManager, error) {
	if cmgr == nil {
		lock.Lock()
		defer lock.Unlock()
		if cmgr == nil {
			_cmgr, err := initConnectionManager()
			if err != nil {
				return nil, err
			}
			cmgr = _cmgr
		}
	}
	return cmgr, nil
}
