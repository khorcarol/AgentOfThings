//go:generate fyne bundle -o bundled.go --package frontend ../assets/Inter_24pt-Bold.ttf
//go:generate fyne bundle -o bundled.go --package frontend --append ../assets/blank-profile.png

package frontend

import (
	// "go/format"

	"image/color"
	"io"
	"log"
	"strconv"

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
	"github.com/khorcarol/AgentOfThings/internal/middle"
)

// Custom colors
var (
	primaryColor    = color.NRGBA{R: 0, G: 120, B: 215, A: 255}   // Blue
	backgroundColor = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Off-white
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

func getImage(interests []api.Interest) *string {
	for i := range interests {
		if interests[i].Image != nil {
			return interests[i].Image
		}
	}
	return nil
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
			profileImage := &canvas.Image{}
			profileImage.SetMinSize(fyne.Size{Width: 200, Height: 200})
			interestImage := &canvas.Image{}
			interestImage.SetMinSize(fyne.Size{Width: 200, Height: 200})

			return container.NewHBox(container.NewVBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Interests: "),
			), profileImage, interestImage)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			friend := currentFriends[i]
			container := o.(*fyne.Container)

			leftBox := container.Objects[0].(*fyne.Container)

			nameLabel := leftBox.Objects[0].(*widget.Label)
			nameLabel.SetText("Name: " + friend.Name)

			interestsLabel := leftBox.Objects[1].(*widget.Label)
			interestsLabel.SetText("Interests: " + formatInterests(friend.User.Interests))

			if friend.Photo.Img != nil {
				image := container.Objects[1].(*canvas.Image)
				image.Image = friend.Photo.Img
				image.FillMode = canvas.ImageFillContain
			}

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
				widget.NewButton("Reject Friend Request", nil),
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

			friendButton := container.Objects[4].(*widget.Button)

			friendButton.OnTapped = func() {
				middle.Seen(user.UserID)
				middle.SendFriendRequest(user.UserID, true)
			}

			rejectButton := container.Objects[5].(*widget.Button)

			rejectButton.OnTapped = func() {
				middle.Seen(user.UserID)
				middle.SendFriendRequest(user.UserID, false)
			}
			rejectButton.Hide()

			if middle.HasIncomingFriendRequest(user.UserID) {
				friendButton.SetText("Accept Friend Request")
				rejectButton.Show()
			} else if middle.HasOutgoingFriendRequest(user.UserID) {
				friendButton.SetText("Friend Request Sent")
				friendButton.Disable()
			}
		},
	)
	return container.NewBorder(nil, nil, nil, nil, usersList)
}

//	func showPopup(win fyne.Window) {
//		dialog.ShowInformation("Notification", "Friend request has been accepted", win)
//	}

var (
	hubsList    *widget.List
	currentHubs []api.Hub
)

func createHubsUI() fyne.CanvasObject {
	hubsList = widget.NewList(
		func() int { return len(currentHubs) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Hub"),
				layout.NewSpacer(),
				widget.NewButton("Open", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// hub := currentHubs[i]
			container := o.(*fyne.Container)

			nameLabel := container.Objects[0].(*widget.Label)
			nameLabel.SetText("Hub " + strconv.Itoa(i+1))

			button := container.Objects[2].(*widget.Button)
			button.OnTapped = func() {
				log.Println("Open hub")
				createHubDialog(currentHubs[i], myWindow)
			}
		},
	)
	return container.NewBorder(nil, nil, nil, nil, hubsList)
}

func createHubDialog(hub api.Hub, myWindow fyne.Window) {

	messages := widget.NewList(
		func() int { return len(hub.Messages) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Message"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			message := hub.Messages[i]
			container := o.(*fyne.Container)

			nameLabel := container.Objects[0].(*widget.Label)
			nameLabel.SetText(message.Contents)

		},
	)

	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Entry", Widget: entry}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			entry.SetText("")
			// TODO: send message to hub ID
		},
	}

	dialog := dialog.NewCustom("Hub", "Close", container.NewVBox(messages, form), myWindow)
	dialog.Resize(fyne.NewSize(500, 200))
	dialog.Show()
}

func onRefreshHubs(hubs []api.Hub) {
	if hubsList != nil {
		currentHubs = hubs
		hubsList.Refresh()
	}
}

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
	middle.Pass(onRefreshFriends, onRefreshUsers, onRefreshHubs)

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
		container.NewTabItem("Hubs", createHubsUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)

}

func Run() {
	onRefreshHubs([]api.Hub{
		{HubID: api.ID{},
			Messages: []api.Message{{Author: api.ID{}, Contents: "Hello"}}},
	})
	myWindow.ShowAndRun()

}
