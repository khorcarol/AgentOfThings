package sources

import (
	"github.com/khorcarol/AgentOfThings/internal/api"
)

func GetInterests() []api.Interest {
	interests := make([]api.Interest, 0)

	// Add new sources below
	interests = append(interests, getPlexInterests()...)


	return interests
}
