package middle

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
)

var (
	hubs = make(map[api.ID]api.Hub)
)

func getHubList() []api.Hub {
	res := []api.Hub{}
	for _, v := range hubs {
		res = append(res, v)
	}
	return res
}

func updateMessages() {
	cmgr := connection.GetCMGR()
	new_hub := <-cmgr.NewMessages

	old_hub, ok := hubs[new_hub.HubID]
	if ok {
		// Hub in hubs
		new_hub.Messages = append(old_hub.Messages, new_hub.Messages...)
		hubs[new_hub.HubID] = new_hub
	} else {
		// New hub
		hubs[new_hub.HubID] = new_hub
	}

	frontend_functions.hubs_refresh(getHubList())
}
