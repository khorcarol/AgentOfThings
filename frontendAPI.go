package main

import "fmt"


// Define struct for interest
type Interest struct {
	Category    string //TODO: enum
	Description string
}

// Define struct for discovered peers (those you find)
type User struct {
	UserID   string
	CommonInterests []Interest
}

// Define struct for peer (connected peers), now includes Discovered struct
type Friends struct {
	user User
	Photo      string
	Name       string
}


//Front calls Middle implements
//stores	
func seen(userID int){

}

func sendFriendRequest(userid UserID){
	
}

//Middle calls Front implemnets

func onRefreshFriends(friends []Friend){
	//refreshes friend UI
}




func onRefreshUsers(users []User ){
	//refreshes user UI
}
