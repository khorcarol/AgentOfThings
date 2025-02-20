package api

import "github.com/khorcarol/AgentOfThings/internal/api/interests"

// A [interests.InterestCategory] and a string to represent the data.
type Interest struct {
	Category    interests.InterestCategory
	Description string
	Image       *string
}

// A [User] is a peer whose interests we have discovered.
type User struct {
	UserID          ID
	CommonInterests []Interest
	Seen            bool
}

// A [Friend] is a [User] who we have requested to be friends with, and who has also requested to be friends with us.
type Friend struct {
	User  User
	Photo string
	Name  string
}

type ID struct {
	Address string
}
