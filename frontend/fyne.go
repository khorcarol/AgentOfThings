//go:generate fyne bundle -o bundled.go --package frontend ../assets/Inter_24pt-Bold.ttf
//go:generate fyne bundle -o bundled.go --package frontend --append ../assets/blank-profile.png

package frontend

import (
	// "go/format"

	"image/color"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
	"github.com/khorcarol/AgentOfThings/internal/middle"
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
	incomingFriendRequests []api.User
)

func frRequest(in []api.User, out []api.User) {
	incomingFriendRequests = in
	outgoingFriendRequests = out
	if incomingFriendsList != nil {
		incomingFriendsList.Refresh()
	}
	outgoingFriendRequests = out
	if outgoingFriendsList != nil {
		outgoingFriendsList.Refresh()
	}
}

func getImage(interests []api.Interest) *string {
	for i := range interests {
		if interests[i].Image != nil {
			return interests[i].Image
		}
	}
	return nil
}

func createFriendRequestsUI() fyne.CanvasObject {
	incomingFriendsList = widget.NewList(
		func() int { return len(incomingFriendRequests) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(
				widget.NewLabel("Anonymous User"),
				widget.NewLabel("Interests: "),
				image,
				widget.NewButton("Accept Friend Request", nil),
				widget.NewButton("Reject Friend Request", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := incomingFriendRequests[i]
			vertContainer := o.(*fyne.Container)

			interests_label := vertContainer.Objects[1].(*widget.Label)
			interests_label.SetText("Interests: " + formatInterests(user.Interests))

			if imageUrl := getImage(user.Interests); imageUrl != nil {
				image := vertContainer.Objects[2].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*imageUrl)
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
			return container.NewVBox(
				widget.NewLabel("Anonymous User"),
				widget.NewLabel("Interests: "),
				image,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := outgoingFriendRequests[i]
			vertContainer := o.(*fyne.Container)

			interests_label := vertContainer.Objects[1].(*widget.Label)
			interests_label.SetText("Interests: " + formatInterests(user.Interests))

			if imageUrl := getImage(user.Interests); imageUrl != nil {
				image := vertContainer.Objects[2].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*imageUrl)
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
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})

			return container.NewVBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Interests: "),
				image,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			friend := currentFriends[i]
			container := o.(*fyne.Container)

			nameLabel := container.Objects[0].(*widget.Label)
			nameLabel.SetText("Name: " + friend.Name)

			interestsLabel := container.Objects[1].(*widget.Label)
			interestsLabel.SetText("Interests: " + formatInterests(friend.User.Interests))

			if imageUrl := getImage(friend.User.Interests); imageUrl != nil {
				image := container.Objects[2].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*imageUrl)
				image.FillMode = canvas.ImageFillContain
			}

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
			return container.NewVBox(
				widget.NewLabel("Anonymous User"),
				layout.NewSpacer(),
				widget.NewLabel("Interests: "),
				image,
				widget.NewButton("Send Friend Request", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := currentUsers[i]
			container := o.(*fyne.Container)

			interests_label := container.Objects[2].(*widget.Label)
			interests_label.Text = "Interests: " + formatInterests(user.Interests)
			interests_label.Refresh()

			if imageUrl := getImage(user.Interests); imageUrl != nil {
				image := container.Objects[3].(*canvas.Image)
				image.Resource, _ = fyne.LoadResourceFromURLString(*imageUrl)
				image.FillMode = canvas.ImageFillContain
			}

			button := container.Objects[4].(*widget.Button)

			button.OnTapped = func() {
				middle.Seen(user.UserID)
				middle.SendFriendRequest(user.UserID, true)
				dialog.ShowInformation("Request Sent", "Friend request sent!", myWindow)
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

var myWindow fyne.Window

func InitLoginForm(callback func(name, interest string, profileImageReader io.ReadCloser)) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	interestsEntry := widget.NewEntry()
	interestsEntry.SetPlaceHolder("What do you like to do?")

	var profileImageReader io.ReadCloser
	profileImageButton := widget.NewButton("Add profile image", func() {
		profileImageDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Printf("Failed to select profile image: %v", err)
			} else {
				profileImageReader = reader
			}
		}, myWindow)

		profileImageDialog.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

		profileImageDialog.Show()
	})

	loginForm := dialog.NewForm(
		"Login",
		"Submit",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Interests", interestsEntry),
			widget.NewFormItem("Profile Image", profileImageButton),
		},
		func(ok bool) {
			if ok {
				callback(nameEntry.Text, interestsEntry.Text, profileImageReader)
			} else {
				myWindow.Close()
			}
		},
		myWindow,
	)

	loginForm.Resize(fyne.NewSize(500, 400))
	loginForm.Show()
}

func Init() {
	middle.Pass(onRefreshFriends, onRefreshUsers, frRequest)

	regularFont := resourceInter24ptBoldTtf

	myApp := app.New()
	myWindow = myApp.NewWindow("Agent of Friends")
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
}

func Run() {
	myWindow.ShowAndRun()
}
