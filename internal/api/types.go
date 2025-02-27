package api

import (
	"github.com/google/uuid"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
)

// A [interests.InterestCategory] and a string to represent the data.
type Interest struct {
	Category    interests.InterestCategory
	Description string
	Image       *string
}

// A [User] is a peer whose interests we have discovered.
type User struct {
	UserID    ID
	Interests []Interest
	Seen      bool
}

// A [Friend] is a [User] who we have requested to be friends with, and who has also requested to be friends with us.
type Friend struct {
	User  User
	Photo string
	Name  string
}

type ID struct {
	Address uuid.UUID
}

// Response to an external friend request
type FriendResponse struct {
	UserID ID // The peer to respond to
	Accept bool
	Data   Friend
}
