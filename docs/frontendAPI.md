

## Data presented

- Person name
- Common interests
    - Films/TV
    - Songs
    - Books
    - etc.

## Structure

### Frontend
- Display people


### Middle
- Store data
- Organise data (rank)

### Backend
- get data
- store personal data



## Data between Middle to Front
### Main Page
- Common interests
    - Catagory
        - Films/TV
        - Music
        - Sports
        - Books
    - Description of interest

### Connected Page
- All connections
    - Pending
        - Same as on main page
    - Connected
        - Photo
        - Name
        - Message

## API Structure

### Structs

#### Interest
- Catagory
- Description


### Functions

#### getCommonInterests
- returns list of Interest

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
