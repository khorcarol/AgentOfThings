//go:generate fyne bundle -o bundled.go --package frontend ../assets/Inter_24pt-Bold.ttf
//go:generate fyne bundle -o bundled.go --package frontend --append ../assets/blank-profile.png

package frontend

import (
	// "time"

	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/middle"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

// Custom colors
var (
	primaryColor    = color.NRGBA{R: 0, G: 120, B: 215, A: 255}   // Blue
	backgroundColor = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Off-white
	listItemColor   = color.White
)

// Custom theme implementation
type customTheme struct {
	fyne.Theme
	regularFont     fyne.Resource
	primaryColor    color.Color
	backgroundColor color.Color
}

func (c *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	if !style.Bold && !style.Italic {
		return c.regularFont
	}
	return c.Theme.Font(style)
}

func (c *customTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return c.primaryColor
	case theme.ColorNameBackground:
		return c.backgroundColor
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

var (
	incomingFriendsList    *widget.List
	outgoingFriendsList    *widget.List
	outgoingFriendRequests []api.User
	IncomingFriendRequests []api.User
)

func frRequest(in []api.User, out []api.User) {
	IncomingFriendRequests = in
	outgoingFriendRequests = out
	if incomingFriendsList != nil {
		incomingFriendsList.Refresh()
	}
	if outgoingFriendsList != nil {
		outgoingFriendsList.Refresh()
	}
}

// func onRefreshIncomingFriendRequests(in []api.User) {

// }

// func onRefreshOutgoingFriendRequests(out []api.User) {

// }

func createFriendRequestsUI() fyne.CanvasObject {
	incomingFriendsList = widget.NewList(
		func() int { return len(IncomingFriendRequests) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(
				widget.NewLabel("User ID"),
				layout.NewSpacer(),
				image,
				widget.NewButton("Accept Friend Request", nil),
				widget.NewButton("Reject Friend Request", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := IncomingFriendRequests[i]
			vertContainer := o.(*fyne.Container)

			if len(user.Interests) > 0 {
				image := vertContainer.Objects[2].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*user.Interests[0].Image)
				image.FillMode = canvas.ImageFillContain
			}

			button := vertContainer.Objects[3].(*widget.Button)
			button.OnTapped = func() {
				log.Println("Accepted friend requests")
				middle.SendFriendRequest(user.UserID, true)
			}

			button2 := vertContainer.Objects[4].(*widget.Button)
			button2.OnTapped = func() {
				log.Println("Rejected friend requests")
				middle.SendFriendRequest(user.UserID, false)
			}
		},
	)

	outgoingFriendsList = widget.NewList(
		func() int { return len(outgoingFriendRequests) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(container.NewHBox(
				widget.NewLabel("User ID"),
				layout.NewSpacer(),
			),
				image,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := outgoingFriendRequests[i]
			vertContainer := o.(*fyne.Container)

			if len(user.Interests) > 0 {
				image := vertContainer.Objects[1].(*canvas.Image)

				image.Resource, _ = fyne.LoadResourceFromURLString(*user.Interests[0].Image)
				image.FillMode = canvas.ImageFillContain
			}
		},
	)
	return container.NewGridWithColumns(2,
		container.NewBorder(
			widget.NewLabel("Incoming Friend Requests"),
			nil, nil, nil, incomingFriendsList,
		),
		container.NewBorder(
			widget.NewLabel("Outgoing Friend Requests"),
			nil, nil, nil, outgoingFriendsList,
		))
}

var (
	friendsList    *widget.List
	usersList      *widget.List
	currentFriends []api.Friend
	currentUsers   []api.User
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
			interestsLabel.SetText(formatInterests(middle.CommonInterests(friend.User.UserID)))
		},
	)

	return container.NewBorder(nil, nil, nil, nil, friendsList)
}

func onRefreshFriends(friends []api.Friend) {
	currentFriends = friends
	if friendsList != nil {
		friendsList.Refresh()
	}
}

func onRefreshUsers(users []api.User) {
	currentUsers = users
	if usersList != nil {
		usersList.Refresh()
	}
}

func formatInterests(interests []api.Interest) string {
	result := ""
	for _, interest := range interests {
		result += interest.Description + "\n"
	}
	return result
}

func createUsersUI(myWindow fyne.Window) fyne.CanvasObject {
	usersList = widget.NewList(
		func() int { return len(currentUsers) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(container.NewHBox(
				widget.NewLabel("User ID"),
				layout.NewSpacer(),
				widget.NewButton("Learn More", nil),
			),
				image,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := currentUsers[i]
			vertContainer := o.(*fyne.Container)
			container := vertContainer.Objects[0].(*fyne.Container)

			button := container.Objects[2].(*widget.Button)
			button.OnTapped = func() {
				middle.Seen(user.UserID)
				showUserDetailsDialog(user, myWindow)
			}

			if len(user.Interests) > 0 {
				image := vertContainer.Objects[1].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*user.Interests[0].Image)
				image.FillMode = canvas.ImageFillContain
			}
		},
	)
	return container.NewBorder(nil, nil, nil, nil, usersList)
}

func showUserDetailsDialog(user api.User, parent fyne.Window) {
	content := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Interests:"),
	)

	for _, interest := range user.Interests {
		log.Println(interest.Description)
		content.Add(widget.NewLabel("- " + interests.String(interest.Category) + ": " + interest.Description))
	}

	var userDetailsDialog dialog.Dialog

	sendBtn := widget.NewButton("Send Friend Request", func() {
		middle.SendFriendRequest(user.UserID, true)
		dialog.ShowInformation("Request Sent", "Friend request sent!", parent)
	})
	closeBtn := widget.NewButton("Close", func() {
		userDetailsDialog.Hide()
	})

	buttons := container.NewHBox(
		sendBtn,
		closeBtn,
	)
	content.Add(buttons)
	userDetailsDialog = dialog.NewCustomWithoutButtons("User Details",
		content,
		parent,
	)

	userDetailsDialog.Show()
}

//	func showPopup(win fyne.Window) {
//		dialog.ShowInformation("Notification", "Friend request has been accepted", win)
//	}
func ShowLoginForm(window fyne.Window) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	interestsEntry := widget.NewEntry()
	interestsEntry.SetPlaceHolder("What do you like to do?")

	loginForm := dialog.NewForm(
		"Login",
		"Submit",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Interests", interestsEntry),
		},
		func(ok bool) {
			if ok {
				personal.AddInterest(api.Interest{Category: 4, Description: interestsEntry.Text})

				var descriptions []string
				for _, interest := range personal.GetSelf().User.Interests {
					descriptions = append(descriptions, interest.Description)
				}

				dialog.ShowInformation("Welcome", "Hello"+"!\nInterests: "+strings.Join(descriptions, ", "), window)
			} else {
				window.Close()
			}
			connection_manager, err := connection.GetCMGR()
			if err != nil {
				log.Fatal("Failed to initialise ConnectionManager:", err)
			}
			middle.Start()
			connection_manager.StartDiscovery()
		},
		window,
	)

	loginForm.Resize(fyne.NewSize(500, 400))
	loginForm.Show()
}

func Main() {
	regularFont := resourceInter24ptBoldTtf

	myApp := app.New()
	myWindow := myApp.NewWindow("Agent of Friends")
	myWindow.Resize(fyne.NewSize(800, 600))
	myApp.Settings().SetTheme(&customTheme{
		Theme:           theme.DefaultTheme(),
		regularFont:     regularFont,
		primaryColor:    primaryColor,
		backgroundColor: backgroundColor,
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Peers", createUsersUI(myWindow)),
		container.NewTabItem("Friends", createFriendsUI()),
		container.NewTabItem("Friend Requests", createFriendRequestsUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	ShowLoginForm(myWindow)
	myWindow.ShowAndRun()
}

func Init() {
	middle.Pass(onRefreshFriends, onRefreshUsers, frRequest)
}
