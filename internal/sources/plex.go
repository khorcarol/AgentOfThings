package sources

import (
	"context"
	"fmt"
	"os"

	"github.com/LukeHagar/plexgo"
	"github.com/khorcarol/AgentOfThings/internal/api"
	. "github.com/khorcarol/AgentOfThings/internal/api/interests"
)

// Fetches session history and writes recently viewed films to a slice of interests.
//
// Requires a local plex instance
func GetPlexInterests() []api.Interest {
	ctx := context.Background()

	s := plexgo.New(
		plexgo.WithServerURL("http://localhost:32400"),
	)

	res, err := s.Sessions.GetSessionHistory(ctx, nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to plex instance")
	}
	if res != nil && res.Object != nil {
		interests := make([]api.Interest, *res.Object.MediaContainer.Size)

		for i, media := range res.Object.MediaContainer.Metadata {
			interests[i] = api.Interest{Category: FilmTV, Description: *media.Title}
		}

		return interests
	}

	return []api.Interest{}
}
