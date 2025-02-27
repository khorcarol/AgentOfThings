package sources

import (
	"context"
	"fmt"
	"os"

	"github.com/LukeHagar/plexgo"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
)

var plexSourceName = "plex"

// Fetches session history and writes recently viewed films to a slice of interests.
//
// Requires a local plex instance
func GetPlexInterests() []api.Interest {
	cached := getCachedSourceInterests(plexSourceName)

	if cached != nil {
		println("Using cached plex data")
		return *cached
	} else {
		println("No cached plex data available")
	}

	ctx := context.Background()

	server_url := "http://localhost:32400"

	s := plexgo.New(
		plexgo.WithServerURL(server_url),
	)

	res, err := s.Sessions.GetSessionHistory(ctx, nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to plex instance")
	}
	if res != nil && res.Object != nil {
		var plex_interests []api.Interest

		interests_contained := map[string]bool{}

		for _, media := range res.Object.MediaContainer.Metadata {
			var image_url *string
			if media.Thumb != nil {
				thumbnail_url := server_url + *media.Thumb
				image_url = &thumbnail_url
			}
			if media.Key != nil && !interests_contained[*media.Key] {
				interests_contained[*media.Key] = true

				plex_interests = append(plex_interests, api.Interest{
					Category:    interests.FilmTV,
					Description: *media.Title,
					Image:       image_url,
				})
			}
		}

		err = cacheSourceInterests(plex_interests, plexSourceName)
		if err != nil {
			println("Failed caching plex data")
			println(err.Error())
		}

		return plex_interests
	}

	return []api.Interest{}
}
