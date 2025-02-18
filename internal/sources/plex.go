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

	server_url := "http://localhost:32400"

	s := plexgo.New(
		plexgo.WithServerURL(server_url),
	)

	res, err := s.Sessions.GetSessionHistory(ctx, nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to plex instance")
	}
	if res != nil && res.Object != nil {
		var interests []api.Interest

		interests_contained := map[string]bool{}

		for _, media := range res.Object.MediaContainer.Metadata {
			var image_url *string
			if media.Thumb != nil {
				thumbnail_url := server_url + *media.Thumb
				image_url = &thumbnail_url
			}
			if media.Key != nil && !interests_contained[*media.Key] {
				interests_contained[*media.Key] = true

				interests = append(interests, api.Interest{
					Category:    FilmTV,
					Description: *media.Title,
					Image:       image_url,
				})
			}
		}

		return interests
	}

	return []api.Interest{}
}
