package main

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Interest struct {
	Category    string
	Description string
}

type User struct {
	UserID          string
	CommonInterests []Interest
}

type Friend struct {
	User  User
	Photo string
	Name  string
}

var (
	friendsList   *widget.List
	usersList     *widget.List
	currentFriends []Friend
	currentUsers   []User
)



// Add these custom colors
var (
	primaryColor    = color.NRGBA{R: 0, G: 120, B: 215, A: 255} // Blue
	backgroundColor = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Off-white
	listItemColor   = color.White
)


func createFriendsUI() fyne.CanvasObject {
	friendsList = widget.NewList(
		func() int { return len(currentFriends) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				canvas.NewImageFromResource(nil),
				container.NewVBox(
					widget.NewLabel("Name"),
					widget.NewLabel("Common Interests"),
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			friend := currentFriends[i]
			container := o.(*fyne.Container)
			img := container.Objects[0].(*canvas.Image)
			details := container.Objects[1].(*fyne.Container)

			img.Resource = nil
			img.SetMinSize(fyne.NewSize(48, 48))

			nameLabel := details.Objects[0].(*widget.Label)
			nameLabel.SetText(friend.Name)
			
			interestsLabel := details.Objects[1].(*widget.Label)
			interestsLabel.SetText(formatInterests(friend.User.CommonInterests))
		},
	)

	return container.NewBorder(nil, nil, nil, nil, friendsList)
}

func onRefreshFriends(friends []Friend) {
	currentFriends = friends
	if friendsList != nil {
		friendsList.Refresh()
	}
}

func onRefreshUsers(users []User) {
	currentUsers = users
	if usersList != nil {
		usersList.Refresh()
	}
}

func formatInterests(interests []Interest) string {
	result := ""
	for _, interest := range interests {
		result += interest.Category + ": " + interest.Description + "\n"
	}
	return result
}

func createUsersUI(myWindow fyne.Window) fyne.CanvasObject {
	usersList = widget.NewList(
		func() int { return len(currentUsers) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("User ID"),
				layout.NewSpacer(),
				widget.NewButton("Learn More", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := currentUsers[i]
			container := o.(*fyne.Container)
			
			idLabel := container.Objects[0].(*widget.Label)
			idLabel.SetText(user.UserID)

			button := container.Objects[2].(*widget.Button)
			button.OnTapped = func() {
				Seen(user.UserID)
				showUserDetailsDialog(user, myWindow)
			}
		},
	)
	return container.NewBorder(nil, nil, nil, nil, usersList)
}

func showUserDetailsDialog(user User, parent fyne.Window) {
	content := container.NewVBox(
		widget.NewLabelWithStyle("User Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("ID: " + user.UserID),
		widget.NewLabel("Common Interests:"),
	)

	for _, interest := range user.CommonInterests {
		content.Add(widget.NewLabel("- " + interest.Category + ": " + interest.Description))
	}

	sendBtn := widget.NewButton("Send Friend Request", func() {
		SendFriendRequest(user.UserID)
		dialog.ShowInformation("Request Sent", "Friend request sent to "+user.UserID, parent)
	})

	content.Add(sendBtn)

	dialog.ShowCustom(
		"",
		"Close",
		content,
		parent,
	)
}

func Seen(userID string) {
	println("Marking user as seen:", userID)
}

func SendFriendRequest(userID string) {
	println("Sending friend request to:", userID)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Social Connect")
	myWindow.Resize(fyne.NewSize(800, 600))

	onRefreshUsers([]User{
		{
			UserID: "user123",
			CommonInterests: []Interest{
				{Category: "Sports", Description: "Basketball"},
				{Category: "Music", Description: "Jazz"},
			},
		},
	})

	onRefreshFriends([]Friend{
		{
        User: User{
            UserID: "123",
            CommonInterests: []Interest{{"Sports", "Basketball"}},
        },
        Name: "John Doe",
        Photo: "path/to/image.jpg",
    },
		
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Peers", createUsersUI(myWindow)),
		container.NewTabItem("Friends", createFriendsUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}


