package api

import (
	"time"

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
	User    User
	Photo   ImageData
	Name    string
	Contact string
}

type ID struct {
	Address uuid.UUID
}

type FriendRequest struct {
	Friend   Friend
	Accepted bool
}

type Hub struct {
	HubName  string
	HubID    ID
	Messages []Message
}

type Message struct {
	Author    ID
	Contents  string
	Timestamp time.Time
}

func (id ID) String() string {
	return id.Address.String()
}

func (id ID) MarshalText() ([]byte, error) {
	return id.Address.MarshalText()
}

func (id *ID) UnmarshalText(text []byte) error {
	return id.Address.UnmarshalText(text)
}
