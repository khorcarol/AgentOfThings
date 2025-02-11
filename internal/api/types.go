package api

import "github.com/khorcarol/AgentOfThings/internal/api/interests"

// Define struct for interest
type Interest struct {
	Category    interests.InterestCategory
	Description string
}

// Define struct for discovered peers (those you find)
type User struct {
	UserID          string
	CommonInterests []Interest
}

// Define struct for peer (connected peers), now includes Discovered struct
type Friends struct {
	user  User
	Photo string
	Name  string
}
